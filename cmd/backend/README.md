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
Build the app with the following command: `go build`