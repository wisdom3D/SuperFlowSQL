import os
import json
from dotenv import load_dotenv

# Charger les variables du fichier .env
load_dotenv()

# --- Récupération des variables ---
PG_HOST = os.getenv("POSTGRES_HOST", "postgres")
PG_PORT = os.getenv("POSTGRES_PORT", "5432")
PG_USER = os.getenv("POSTGRES_USER", "airflow")
PG_PASS = os.getenv("POSTGRES_PASSWORD", "airflow")
PG_DB = os.getenv("POSTGRES_DB_AIRFLOW", "airflow")

# --- Définition du contenu JSON ---
servers_data = {
  "Servers": {
    "1": {
      "Name": "SuperFlowSQL - PostgreSQL",
      "Group": "SuperFlowSQL Stack",
      "Port": int(PG_PORT),
      "Username": PG_USER,
      "Host": PG_HOST,
      "SSLMode": "disable",
      "MaintenanceDB": PG_DB,
      # Le mot de passe n'est pas nécessaire ici, car l'utilisateur devra le fournir
      # lors de la première connexion via l'interface PgAdmin.
    }
  }
}

# --- Écriture du fichier ---
# Nous allons écrire ce fichier dans le dossier 'config' ou directement à la racine
# selon la structure que vous préférez. Je vais le mettre à la racine pour la simplicité
# et le nommer pgadmin_servers.json comme précédemment.
output_file = "pgadmin_servers.json"

try:
    with open(output_file, 'w') as f:
        json.dump(servers_data, f, indent=4)
    print(f"✅ Fichier de configuration PgAdmin généré : {output_file}")
except Exception as e:
    print(f"❌ Erreur lors de la génération du fichier : {e}")