# Todo items
- Community details page: we need to load the community details, as well as the posts

## pending endpoints
- getCommunityMetadata
- getCommunityPosts (returns list of posts - title, shortenedContent, username, date, vote, commentCount)
- getPostComments
- getTimelinePosts


## Deployment
# Build the image
docker build -t threads99 .

# Run the container
docker run -d -p 8600:8090 threads99

docker build -t ysodiqakanni/swizzle-api:latest .
docker push ysodiqakanni/swizzle-api:latest

Kubernetes:
kubectl apply -f .\quibbles-api-deployment.yaml
Force restart deployment to pick new image: `kubectl rollout restart deployment/quibbles-api`
Exec into pod:
kubectl exec -i -t -n default quibble-api-64585c6467-gtw8l -c quibbles-api -- sh -c "clear; (bash || ash || sh)"

Healthcheck:
install curl (for Alpine): `apk add --no-cache curl`
curl http://localhost:8090/api/healthcheck


In the server,
docker pull ysodiqakanni/swizzle-api:latest
docker run -d -p 8600:8090 --name swizzle-api ysodiqakanni/swizzle-api:latest

press20Five
#mongosh mongodb+srv://root:Password@helloworldcluster.zndnutk.mongodb.net/test