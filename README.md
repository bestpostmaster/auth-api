# Auth API

API Go/MySQL exposant `POST /api/login` et `POST /api/add-user`. Le champ JSON `username` correspond à la colonne `user.email` du schéma existant. Le compte doit être actif et son mot de passe doit être un hash bcrypt.

La migration `migrations/000001_create_users_table.up.sql` crée la table `user` ainsi qu'un index unique sur `email`, nécessaire pour garantir qu'un identifiant de connexion ne désigne qu'un seul compte.

## Configuration

Copier `.env.example` vers `.env` et adapter ses valeurs. Le fichier `.env` est chargé automatiquement au démarrage ; une variable déjà exportée dans le shell reste prioritaire. Générer une clé RSA locale (ne pas la committer) :

```bash
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out private.pem
openssl rsa -pubout -in private.pem -out public.pem
```

Les variables obligatoires sont `MYSQL_DSN` et soit `JWT_PRIVATE_KEY_FILE`, soit `JWT_PRIVATE_KEY` (contenu PEM, les `\\n` littéraux sont acceptés). La durée d'un JWT est de 60 minutes par défaut.

## Migrations de la base de données

Installer la CLI `golang-migrate` avec le support MySQL :

```bash
go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

La CLI attend une URL commençant par `mysql://`, contrairement à la variable `MYSQL_DSN` utilisée par l'API. Définir l'URL dans le terminal en l'entourant d'apostrophes, notamment lorsque le mot de passe contient un caractère `$` :

```bash
export MIGRATIONS_DATABASE_URL='mysql://user:mot-de-passe@tcp(127.0.0.1:3306)/auth'
```

Appliquer toutes les migrations disponibles :

```bash
migrate -path migrations -database "$MIGRATIONS_DATABASE_URL" up
```

Afficher la version actuellement appliquée :

```bash
migrate -path migrations -database "$MIGRATIONS_DATABASE_URL" version
```

Annuler la dernière migration :

```bash
migrate -path migrations -database "$MIGRATIONS_DATABASE_URL" down 1
```

Créer une nouvelle paire de fichiers de migration :

```bash
migrate create -ext sql -dir migrations -seq nom_de_la_migration
```

Chaque migration génère un fichier `.up.sql` pour l'application et un fichier `.down.sql` pour le rollback. La commande `down` modifie la base de données et doit être utilisée avec précaution.

## Démarrage

```bash
go run .
```

La réponse d'une authentification réussie est :

```json
{"token":"<JWT RS256>","userId":2}
```

Des identifiants incorrects ou un compte inactif produisent une réponse HTTP `401`. Le navigateur provenant de `http://localhost:4201` est autorisé par CORS.

## Création d'un utilisateur

```http
POST /api/add-user
Content-Type: application/json

{"username":"test","password":"abcd12345"}
```

La réponse est `201 Created` avec `{"userId":<id>}`. Le mot de passe doit contenir entre 8 et 72 caractères. Il est stocké sous forme de hash bcrypt. Un identifiant existant produit `409 Conflict`. Conformément au schéma, un utilisateur nouvellement créé est inactif (`is_active = 0`).
