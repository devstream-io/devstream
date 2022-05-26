kind: Service
apiVersion: v1
metadata:
  name: [[.AppName]]
spec:
  selector:
    app: [[.AppName]]
  ports:
    - port: 80
      targetPort: 8080
