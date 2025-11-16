# 🚀 Data Platform - SuperFlowSQL

Une plateforme de data orchestration complète basée sur **Apache Airflow**, **PostgreSQL**, **PgAdmin** et **Apache Superset**. Cette solution permet d'automatiser vos pipelines de données et de les monitorer en temps réel.

## 📋 Architecture

```
┌──────────────────────────────────────────────────┐
│          Data Platform Architecture              │
├──────────────────────────────────────────────────┤
│                                                  │
│  🌐 Airflow Webserver (Port 8080)               │
│     └─ Orchestration & Monitoring               │
│                                                  │
│  ⚙️  Airflow Scheduler                          │
│     └─ DAGs Execution                           │
│                                                  │
│  📊 PostgreSQL (Port 5432)                      │
│     └─ Data Storage & Airflow Metadata          │
│                                                  │
│  🔧 PgAdmin (Port 5050)                         │
│     └─ Database Management UI                   │
│                                                  │
│  📈 Apache Superset (Port 8088)                 │
│     └─ Data Visualization & Analytics           │
│                                                  │
└──────────────────────────────────────────────────┘
```

## 🛠️ Prérequis

- **Docker** (version 20.10+)
- **Docker Compose** (version 1.29+)
- **Git** (pour cloner le projet)
- Au minimum **2GB de RAM disponible**

### Installation sur Windows

1. **Installer Docker Desktop** : https://www.docker.com/products/docker-desktop
2. **Installer Git** : https://git-scm.com/download/win

## 🚀 Démarrage rapide

### 1️⃣ Configuration initiale

Clonez ou naviguez vers le répertoire du projet :

```bash
cd data-platform
```

### 2️⃣ modifier le fichier `.env`


**⚠️ Note de sécurité** : Modifiez les mots de passe par défaut en production !

### 3️⃣ Générer la configuration PgAdmin

Avant de lancer les conteneurs, générez le fichier de configuration PgAdmin :

```bash
python generate_pgadmin_config.py
```

Ce script crée un fichier `pgadmin_servers.json` qui configure automatiquement la connexion PostgreSQL.

### 4️⃣ Lancer la plateforme

Lancez tous les services avec une seule commande :

```bash
docker-compose up --build -d
```

**Détails de la commande :**
- `up` : Crée et démarre les conteneurs
- `--build` : Reconstruit les images Docker
- `-d` : Mode détaché (exécution en arrière-plan)

### 5️⃣ Vérifier le statut

Consultez l'état des conteneurs :

```bash
docker-compose ps
```

Attendez que tous les services soient en état `healthy` ou `running` (cela peut prendre 1-2 minutes).

## 🌐 Accès aux services

Une fois les conteneurs lancés, accédez aux services :
voir les identifiants dans le fichier .env

| Service | URL | Identifiants |
|---------|-----|--------------|
| **Airflow Webserver** | http://localhost:8080 | `AIRFLOW_USER` / `AIRFLOW_PASSWORD` |
| **PgAdmin** | http://localhost:5050 | `PGADMIN_EMAIL` / `PGADMIN_PASSWORD` |
| **PostgreSQL** | localhost:5432 | `POSTGRES_USER` / `POSTGRES_PASSWORD` |
| **Superset** | http://localhost:8088 | `SUPERSET_USER` / `SUPERSET_PASSWORD`|


### Arrêter la plateforme

```bash
docker-compose down
```

## 🔍 Accéder à la base de données

### Via PgAdmin (Interface Web)
1. Accédez à http://localhost:5050
2. Authentifiez-vous
3. La connexion PostgreSQL devrait être préconfigurée

### Via ligne de commande
```bash
docker-compose exec postgres psql -U airflow -d airflow
```

Commandes SQL utiles :
```sql
-- Lister les bases de données
\l

-- Se connecter à une base
\c airflow

-- Lister les tables
\dt

-- Exécuter une requête
SELECT * FROM users;
```

## 🐛 Troubleshooting

### Les conteneurs ne démarrent pas
```bash
# Vérifiez les logs
docker-compose logs

# Nettoyez et réessayez
docker-compose down -v
docker-compose up --build -d
```

### Port déjà utilisé
Si un port est occupé (ex: 8080), modifiez le fichier `.env` :
```env
AIRFLOW_PORT=8081  # Changez 8080 en 8081
```

### Erreur de connexion PostgreSQL
```bash
# Vérifiez que PostgreSQL est actif
docker-compose ps postgres

# Redémarrez PostgreSQL
docker-compose restart postgres
```

### Les DAGs n'apparaissent pas dans Airflow
1. Vérifiez que les fichiers sont dans `dags/`
2. Vérifiez qu'ils respectent la syntaxe Airflow
3. Consultez les logs du webserver : `docker-compose logs airflow-webserver`

## 📊 Configuration Superset (optionnel)

Pour visualiser vos données :

Ajouter PostgreSQL comme source de données dans Superset :

Host = postgres nom du service dans docker-compose.yml

Port = 5432 

database = airflow

Username = airflow 

password = airflow_password (ou celui défini dans .env)

1. Accédez à http://localhost:8088
2. Créez un compte administrateur
3. Ajoutez PostgreSQL comme source de données
4. Créez des dashboards

## 🔐 Sécurité

**⚠️ IMPORTANT pour la production :**

1. Modifiez TOUS les mots de passe par défaut
2. Utilisez une clé Fernet robuste pour Airflow
3. Activez HTTPS sur PgAdmin et Superset
4. Restreignez les accès réseau (pare-feu)
5. Mettez en place une sauvegarde régulière des données

## 📝 Génération de clé Fernet pour Airflow

Si vous devez générer une nouvelle clé :

```bash
python -c "from cryptography.fernet import Fernet; print(Fernet.generate_key().decode())"
```

## 🤝 Support et Documentation

- **Apache Airflow** : https://airflow.apache.org/docs/
- **PostgreSQL** : https://www.postgresql.org/docs/
- **PgAdmin** : https://www.pgadmin.org/docs/
- **Apache Superset** : https://superset.apache.org/docs/6.0.0/intro
---

**❤️ SuperFlowSQL**
