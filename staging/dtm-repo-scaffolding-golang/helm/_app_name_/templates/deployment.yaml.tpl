apiVersion: apps/v1
kind: Deployment
metadata:
  name: [[.AppName]]
  labels:
    app: [[.AppName]]
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: [[.AppName]]
  template:
    metadata:
      labels:
        app: [[.AppName]]
    spec:
      containers:
      - name: [[.AppName]]
        image: "[[.ImageRepo]]:{{ .Chart.AppVersion }}"
        ports:
        - containerPort: 8080
