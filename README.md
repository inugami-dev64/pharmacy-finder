# Pharmacy finder

Pharmacy finder helps Estonian trans and non-binary people to find pharmacies that have a good track record of accepting prescriptions from [Imago](https://www.imago.tg) and [GenderGP](https://www.gendergp.com).

## Build and deployment

Pharmacy finder is best used with docker. Before building the image, you will need to create keys for [reCaptcha v2](https://developers.google.com/recaptcha/intro). Due to the nature of the project, the sitekey will get built into the executable and thus, it needs to be passed to docker build as a build argument.

```bash
$ docker build --build-arg RECAPTCHA_SITE_KEY=<mykey> -t pharmafinder .
```

When running the container, the server listens on port `8080`. Additionally you will need to pass environment variables into your docker container (see [deploy/.env.sample](deploy/.env.sample) for more information)