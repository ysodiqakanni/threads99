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

In the server,
docker pull ysodiqakanni/swizzle-api:latest
docker run -d -p 8600:8090 --name swizzle-api ysodiqakanni/swizzle-api:latest

press20Five