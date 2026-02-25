#--------------------------------------
# Exemple de DAG Airflow : Modèle Générique
#-------------------------------------

from __future__ import annotations

import pendulum

from airflow.models.dag import DAG
from airflow.operators.bash import BashOperator
from airflow.operators.python import PythonOperator
# DummyOperator est souvent utilisé pour des points de départ/fin ou des regroupements
# from airflow.operators.dummy import DummyOperator 

# --- 1. Définition des arguments par défaut ---
# Ces arguments sont passés à tous les opérateurs du DAG.
default_args = {
    "owner": "airflow",
    "depends_on_past": False,
    "email": ["mon_email@example.com"],
    "email_on_failure": False,
    "email_on_retry": False,
    "retries": 1,
    # timedelta(minutes=5) : attente de 5 minutes avant une nouvelle tentative
    # "retry_delay": timedelta(minutes=5),
}

# --- 2. Définition de la fonction Python pour le PythonOperator ---
def ma_fonction_a_executer(execution_date, **kwargs):
    """
    Exemple de fonction Python qui sera exécutée par une tâche.
    Elle affiche la date d'exécution.
    """
    print(f"La tâche Python s'exécute à la date logique: {execution_date}")
    # Vous pouvez accéder à d'autres informations du contexte Airflow via **kwargs
    # Par exemple, pour l'ID du DAG : kwargs['dag_run'].dag_id
    return "Résultat de la fonction"


# --- 3. Définition du DAG ---
with DAG(
    dag_id="dag_generique_modele",
    start_date=pendulum.datetime(2023, 1, 1, tz="UTC"),
    # Le DAG s'exécutera tous les jours à minuit
    schedule="0 0 * * *",
    # Ou 'timedelta(days=1)' si vous n'utilisez pas de cron
    catchup=False,  # Important : si False, n'exécute pas les runs manqués avant start_date
    default_args=default_args,
    tags=["exemple", "generique"],
    doc_md=__doc__, # Utilisez la docstring du fichier pour la documentation du DAG
) as dag:
    # --- 4. Définition des Tâches (Opérateurs) ---

    # Tâche 1: Démarrer (souvent un DummyOperator, ici un simple Bash)
    tache_demarrage = BashOperator(
        task_id="demarrer_le_pipeline",
        bash_command='echo "Démarrage du pipeline!"',
    )

    # Tâche 2: Exécuter une commande Bash
    tache_traitement_bash = BashOperator(
        task_id="traitement_donnees_bash",
        # {{ ds }} est une macro Airflow pour la date d'exécution
        bash_command='echo "Traitement Bash pour la date: {{ ds }}" && sleep 5',
    )
    
    # Tâche 3: Exécuter une fonction Python
    tache_traitement_python = PythonOperator(
        task_id="traitement_donnees_python",
        python_callable=ma_fonction_a_executer,
        # Vous pouvez passer des arguments à la fonction si nécessaire :
        # op_kwargs={'argument_key': 'argument_value'},
    )

    # Tâche 4: Fin du pipeline
    tache_fin = BashOperator(
        task_id="terminer_le_pipeline",
        bash_command='echo "Pipeline terminé avec succès!"',
    )


    # --- 5. Définition des dépendances (Ordre d'exécution) ---
    
    # Séquentiel: T1 -> T2 -> T3 -> T4
    # tache_demarrage >> tache_traitement_bash >> tache_traitement_python >> tache_fin

    # Ordre plus complexe :
    # Tâche de Démarrage doit se terminer avant les deux tâches de traitement
    tache_demarrage >> [tache_traitement_bash, tache_traitement_python]
    
    # Les deux tâches de traitement doivent se terminer avant la Tâche de Fin
    [tache_traitement_bash, tache_traitement_python] >> tache_fin

    # L'ordre résultant est :
    # tache_demarrage 
    #        |
    #        --> tache_traitement_bash
    #        --> tache_traitement_python
    #                    |
    #                    --> tache_fin