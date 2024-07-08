# venmo-calculator

This is a simple app that helps you split a bill with your friends (presumably via venmo request). This app is meant to be self-hosted and we provide a docker compose stack to make deployment simple.

In keeping with the simplicity, no persistent storage is supported or required.

Proportionality is achieved by calculating the proportion of the pre tax, tip, or fee subtotal attributable to each person, and then distributing the final total (post-tax, tip, or fee) using that calculation.

## Architecture/Technologies
- This app is designed to deploy via docker-compose.
- The backend is written in `golang` and uses the `chi` router - [(repo here)](https://github.com/go-chi/chi)
- The frontend is a simple SPA written in javascript using the [vue.js](https://vuejs.org/) framework and built using [Vite](https://vitejs.dev/) and styled using vanilla CSS in component scoped styles (thanks, vue!)

## Deployment
When you decide what the URL for your instance of this app will be, [modify this file](./web/.env.production) with the appropriate base URL.

You'll also need to change the `WEB_URL` environment variable in the [docker-compose.yml](docker-compose.yml).

Then - simply start the docker compose stack:

```
docker compose up -d
```

## Development
- The frontend comes with its own dev instructions, which require `npm`. I recommend managing your `npm` versions using [nvm](https://github.com/nvm-sh/nvm)
- The backend is an uncomplicated `golang` app - I've historically managed `golang` versions using [gvm](https://github.com/moovweb/gvm), but for this product I tried using [g](https://github.com/stefanmaric/g) and enjoyed it.
