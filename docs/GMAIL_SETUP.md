# Gmail Credentials Setup

drfin reads bank notifications through the Gmail API using an OAuth 2.0 client.
This guide explains how to create the client credentials JSON
(`resources/drfin-secret.json`).

> Never commit `resources/drfin-secret.json` or `token.json`. Both must stay in `.gitignore`.

## 1. Create or select a Google Cloud project

1. Open <https://console.cloud.google.com/>.
2. Use the project picker at the top, then **New Project**.
3. Name it (for example `drfin`) and click **Create**, then select it.

## 2. Enable the Gmail API

1. Go to **APIs & Services → Library**.
2. Search for **Gmail API** and click **Enable**.

## 3. Configure the OAuth consent screen

1. Go to **APIs & Services → OAuth consent screen** (may appear as **Google Auth Platform**).
2. Choose user type **External** and fill in the app name and your email.
3. Add the scope the app uses: <!-- VERIFY against auth.go --> `https://www.googleapis.com/auth/gmail.readonly`
4. Under **Test users**, add the Gmail address that receives the bank notifications.

> While the app is in **Testing** status, Google expires refresh tokens after
> 7 days, so you will need to re-authorize periodically (delete `token.json` and sign in again).

## 4. Create the OAuth client and download the JSON

1. Go to **APIs & Services → Credentials**.
2. Click **Create Credentials → OAuth client ID**.
3. Application type: <!-- VERIFY against auth.go --> **Desktop app** (or **Web application** if the code uses a redirect URL; in that case add the exact redirect URI the code uses).
4. Name it and click **Create**.
5. Click **Download JSON** in the confirmation dialog (or the download icon next to the client in the Credentials list).
6. Save the file as `resources/drfin-secret.json`.

## 5. First run and token.json

On the first Gmail sync, the app asks you to authorize access and then
creates `token.json`. Run `go run .` from the repository root, since the
credential paths are relative to the working directory.

## Troubleshooting

| Error | Cause |
|---|---|
| `open resources/drfin-secret.json: cannot find the file` | File missing, or app started outside the repo root |
| `access_denied` / app not verified | Your account isn't listed as a test user |
| Token stops working after about a week | Testing-status refresh token expired; delete `token.json` and re-authorize |

## Rotating leaked credentials

1. In **Credentials**, delete the compromised OAuth client.
2. Create a new one (step 4) and replace `resources/drfin-secret.json`.
3. Delete `token.json` and re-authorize.