#
# Final image suitable for deployment. We just take the binary we generated
# from a previous step in the build and copy it into this image.
#
# Required Build Args:
#
#  - APP_BINARY. Path to binary to copy in as the application binary.
#                This path must be relative to the root of the project. For
#                example, dist/myservice.
#
FROM alpine:latest
ARG APP_BINARY

#
# If your final deploy image needs to modified in some other way, do that here.
# For example, adding debugging tools, etc. Be smart in how you structure
# things to maximize layer reuse/caching.
#

COPY ${APP_BINARY} /app
RUN chmod +x /app
CMD ["./app"]
