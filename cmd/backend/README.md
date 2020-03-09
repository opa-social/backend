# Backend

## Getting Started
Run the following commands to get your application credentials. You will need
the `gcloud` utility and have initialized it before.
```bash
$ gcloud iam service-accounts list
$ # Find the Firebase Admin SDK service account in this list and copy the email.
$ gcloud iam service accounts keys create google-services.json \
    --iam-account <service-account@google.com>
```

Next, create the Firebase config:
```bash
cat <<EOF > firebase-config.json
{
    "projectId": "opa-268409",
    "databaseURL": "https://opa-268409.firebaseio.com/"
}
EOF
```

Finally, export the following evironment variables to run the backend.
```bash
$ export GOOGLE_APPLICATION_CREDENTIALS=google-services.json
$ export FIREBASE_CONFIG=firebase-config.json
```

## Building and Running the Application
Build and run the app with the following command: `go run .` 
This will run the application at address `0.0.0.0:9000`.

## Testing the Application
To test the web frontend you will need the following tools:
* [HTTPie]()
* [HTTPie JWT plugin]()
* [jq]()

You will also need to to use the [get-token](../../tools/get-token) tool in
order to get an authentication token.

Follow these steps to get an authentication token to test with.
1. `cd` into `../../tools/get-token`.
2. Go to [Firebase Console General Settings](https://console.firebase.google.com/project/_/settings/general/)
    and copy the "Web API Key".
3. Run `token.sh` in the following manner:
```bash
$ export JWT_AUTH_TOKEN=$(./token.sh --username=<user@email.com> \
    --password=<password> \
    --token=<web-api-token>)
```

Now you can test the application with HTTPie.
```bash
$ http --auth-type=jwt localhost:9000
```