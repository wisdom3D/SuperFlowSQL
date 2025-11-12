from __future__ import annotations

import pendulum
import pandas as pd
import numpy as np

from airflow.decorators import dag, task
from airflow.providers.postgres.hooks.postgres import PostgresHook
from sqlalchemy import create_engine
import os # Nécessaire pour les variables d'environnement si vous utilisez un moteur SQLAlchemy direct

# --- Configuration du DAG ---
TABLE_NAME = "random_users_data"
POSTGRES_CONN_ID = "postgres_default"  # Doit correspondre à votre connexion Airflow

@dag(
    dag_id="data_generation_pandas_to_postgres",
    start_date=pendulum.datetime(2023, 1, 1, tz="UTC"),
    schedule=None,  # Exécution manuelle
    catchup=False,
    tags=["data_ingestion", "pandas", "postgres"],
)
def pipeline_generer_charger_donnees():
    """
    Ce pipeline génère des données aléatoires, les charge dans un DataFrame Pandas,
    puis les insère dans une table PostgreSQL.
    """

    @task
    def generer_et_charger_donnees():
        """
        Génère 100 lignes de données factices, crée un DataFrame Pandas,
        puis utilise le PostgresHook d'Airflow pour charger les données.
        """
        print("--- 1. Génération des données aléatoires ---")
        
        # Nombre de lignes à générer
        N_ROWS = 100
        
        # Création du dictionnaire de données
        data = {
            'user_id': np.arange(1, N_ROWS + 1),
            'age': np.random.randint(20, 60, size=N_ROWS),
            'city': np.random.choice(['Paris', 'Lyon', 'Marseille', 'Nice', 'Lille'], size=N_ROWS),
            'salary': np.round(np.random.normal(50000, 15000, N_ROWS), 2)
        }
        
        # Création du DataFrame Pandas
        df = pd.DataFrame(data)
        print(f"DataFrame créé avec {len(df)} lignes et les colonnes: {df.columns.tolist()}")
        print("\nAperçu des données générées :")
        print(df.head())

        # --- 2. Connexion à PostgreSQL via l'Airflow Hook ---
        
        # Le Hook permet d'interagir avec la base de données en utilisant la connexion Airflow
        pg_hook = PostgresHook(postgres_conn_id=POSTGRES_CONN_ID)
        
        # Utilisation d'un moteur SQLAlchemy pour permettre la fonction to_sql de Pandas
        # Nous récupérons les informations de connexion du Hook pour créer un moteur
        
        # NOTE IMPORTANTE: Pour des raisons de sécurité et de simplicité en production, 
        # il est **fortement recommandé** d'utiliser le Hook Airflow.
        # Cependant, Pandas `to_sql` nécessite un moteur SQLAlchemy.
        
        # Si votre connexion Airflow est simple (comme 'postgresql://user:pass@host:port/db'),
        # vous pouvez essayer de l'obtenir directement :
        # conn_uri = pg_hook.get_uri() 
        # engine = create_engine(conn_uri)
        
        # Sinon, récupérons les composants de la connexion pour construire l'URI
        conn = pg_hook.get_connection(POSTGRES_CONN_ID)
        
        # Assurez-vous que le dialecte est correct (e.g., 'postgresql', 'mysql', etc.)
        conn_uri = f"postgresql+psycopg2://{conn.login}:{conn.password}@{conn.host}:{conn.port}/{conn.schema}"
        engine = create_engine(conn_uri)
        
        print(f"\n--- 3. Chargement dans PostgreSQL dans la table : {TABLE_NAME} ---")

        # Chargement du DataFrame vers la base de données
        df.to_sql(
            name=TABLE_NAME, 
            con=engine, 
            if_exists='replace', # Options : 'fail', 'replace', 'append'
            index=False           # Ne pas inclure l'index du DataFrame comme colonne
        )
        
        print("✅ Chargement terminé avec succès.")
        
    # --- Définition du Workflow (les dépendances) ---
    
    # Dans ce cas simple, il n'y a qu'une seule tâche :
    generer_et_charger_donnees()
    
# Instanciation du DAG
pipeline_generer_charger_donnees()

# --- Requête SQL de Vérification ---
# Si vous voulez une tâche de vérification après le chargement :
"""
@task
def verifier_chargement(table_name):
    pg_hook = PostgresHook(postgres_conn_id=POSTGRES_CONN_ID)
    sql_check = f"SELECT COUNT(*) FROM {table_name};"
    record_count = pg_hook.get_first(sql_check)[0]
    print(f"La table {table_name} contient maintenant {record_count} lignes.")
    assert record_count > 0, "Aucune donnée n'a été insérée !"

verification = verifier_chargement(table_name=TABLE_NAME)
generer_et_charger_donnees() >> verification
"""