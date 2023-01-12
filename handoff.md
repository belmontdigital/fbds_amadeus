# Amadeus Master Schedule/Room Cover Screens

## Work Completed
- API Client integration with near complete support (~90%) for the public documentation of the AHWS Events API
- Caching any data received from the API for 15 minutes
    - Auth token also expires every 15 minutes using refresh token until it is expired (72 hours) and then it goes through normal authenticate method
    - Caching is done in a thread safe manner (using `sync.Map` under the hood), and is saved to disk when the process finishes and loaded when process starts.
- HTML pages rendering using Go Templates
    - Date/Time on screens is rendered using JavaScript so it is always up to date with the client.
- Updated CPanel instance as several packages were years out of date
    - Update was not working without disabling the OwnCloud repository on the server.
- Deployed within CPanel, serving through Apache using Phusion Passenger as a GLS app.
    - Configuration on CPanel side can be found at `/etc/apache2/conf.d/userdata/ssl/2_4/fbapi/fbds-api.belmont.digital` and `/etc/apache2/conf.d/userdata/std/2_4/fbapi/fbds-api.belmont.digital`
    - Application can be found at `/home/fbapi/www/fbds_api`
    - Application can be rebuilt using the `Go` installation in `/home/fbapi/go/bin/go`

## Suggested Follow Up

There are some remaining follow up items that would be good to be resolved post-closure of this project:

1. Time currently defaults to EST since we don't have any TZ information in the timestamps, and can not infer this directly via API without some translation on the Location data.
2. The project was developed in a procedural manner in order to simplify building out the integration. As a follow up, it would be a good idea to separate out the individual concerns of the application into separate packages.
2. API Client currently does not have automated test coverage. Would suggest adding unit tests at a minimum.
3. The logging coverage is pretty good, but currently the logs are all landing in `/var/log/apache2/error.log`, would be good to have this apps logs going to it's own file.
4. (Nice to have, especially if the project grows) - Additional observability improvements could be made such as implementing:
  a. Metrics to track resource utilization, request/response timing buckets, cache hit/miss ratio, etc. to be able to reason about the health of the application as well as areas for improvement.
  b. Tracing
  c. Structured Logging
4. Ideally, use redis for caching, if scaling from a single instance to multiple instances.