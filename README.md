# Auth API

API Go/MySQL exposant `POST /api/login`. Le champ JSON `username` correspond à la colonne `user.email` du schéma existant. Le compte doit être actif et son mot de passe doit être un hash bcrypt.

Le fichier `database-schema/schema.sql` crée également un index unique sur `email`, nécessaire pour garantir qu'un identifiant de connexion ne désigne qu'un seul compte.

## Configuration

Copier `.env.example` vers `.env` et adapter ses valeurs. Le fichier `.env` est chargé automatiquement au démarrage ; une variable déjà exportée dans le shell reste prioritaire. Générer une clé RSA locale (ne pas la committer) :

```bash
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private.pem
openssl rsa -pubout -in private.pem -out public.pem
```

Les variables obligatoires sont `MYSQL_DSN` et soit `JWT_PRIVATE_KEY_FILE`, soit `JWT_PRIVATE_KEY` (contenu PEM, les `\\n` littéraux sont acceptés). La durée d'un JWT est de 60 minutes par défaut.

## Démarrage

```bash
go run .
```

La réponse d'une authentification réussie est :

```json
{"token":"<JWT RS256>","userId":2}
```

Des identifiants incorrects ou un compte inactif produisent une réponse HTTP `401`. Le navigateur provenant de `http://localhost:4201` est autorisé par CORS.
