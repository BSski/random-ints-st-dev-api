<p align="center">
  <h1 align="center">Random Integers' Standard Deviation API</h2>
</p>


<div markdown="1" align="center">

[![Build Status](https://bsski.semaphoreci.com/badges/RandomIntsStDevAPI/branches/main.svg?style=shields&key=b033d465-f52b-4169-8a0c-025b0586f456)](https://random-ints-st-dev-api.herokuapp.com/random/mean?requests=2&length=3)
[![Maintainability](https://api.codeclimate.com/v1/badges/9fd1a96061671d7f7766/maintainability)](https://codeclimate.com/github/BSski/RandomIntsStDevAPI/maintainability)
[![CodeFactor](https://www.codefactor.io/repository/github/bsski/random-ints-st-dev-api/badge)](https://www.codefactor.io/repository/github/bsski/random-ints-st-dev-api)
<!-- [![Demo Uptime](https://img.shields.io/uptimerobot/ratio/7/m792386113-957fa47111c7ccefd3d8617e)](https://random-ints-st-dev-api.herokuapp.com/random/mean?requests=2&length=2) -->
<!-- [![Heroku](https://pyheroku-badge.herokuapp.com/?app=django-music-bsski&style=flat)](https://random-ints-st-dev-api.herokuapp.com/random/mean?requests=2&length=3) -->


<div markdown="1" align="center">    

![Demo Screenshot](https://i.imgur.com/OgMTzZX.png)

</div>


<h3 align="center">
  :fireworks: Live demo :fireworks:
</h3>

<p align="center">
  <b>[Heroku retracted its free tier, therefore all my demos are unavailable until further notice.]</b>
</p>

<!--
<p align="center">
  https://random-ints-st-dev-api.herokuapp.com/random/mean?requests=2&length=3
</p>
-->

<p align="left">

## Table of contents
* [Project description](#scroll-project-description)
* [Technologies used](#hammer-technologies-used)
* [Deployment](#hammer_and_wrench-deployment)
* [Environment variables](#closed_lock_with_key-environment-variables)
* [Room for improvement](#arrow_up-room-for-improvement)
* [Author](#construction_worker-author)


## :scroll: Project description
This is a recruitment task for a Junior Software Developer position.

The project is a REST service supporting the following GET operation:

```/random/mean?requests={r}&length={l}```

which performs `{r}` concurrent requests to random.org API asking for `{l}` number of random integers.

Additionally, the application calculates standard deviation of each drawn integers set and of the sum of all sets.

The project also has a Semaphore CI/CD pipeline with deployment to Heroku via Docker container:
![CI/CD screenshot](https://i.imgur.com/5v0Xufr.png)

I decided to use `runtime.GOMAXPROCS(1)`, since random.org guidelines (https://www.random.org/clients/) prohibit sending simultaneous requests by automated clients. In accordance with the guidelines the timeout is set to 180 seconds.


## :hammer: Technologies used
- random.org API
- Go 1.18
- go-chi/chi/v5 5.0.7
- google/go-cmp 0.5.8
- joho/godotenv 1.4.0
- montanaflynn/stats 0.6.6
- cosmtrek/air 1.40.4
- Docker


## :hammer_and_wrench: Deployment

If you want to run the application, you have to supply your own API key and your email address (the guidelines require random.org API client's email supplied in the "User-Agent" header) in the .env file .

1. Create an `.env` file basing on `.env_sample_file` from the repository. Set `PORT` to 8080.

2. Run `docker run --env-file .env -p 8080:8080 bsski/random-ints-st-dev-api:latest` in the `.env` file directory.

3. Access `localhost:8080/random/mean?requests=2&length=3`. 


## :closed_lock_with_key: Environment variables

To run this project, you have to set up the following environment variables in the `.env` file (**the values below are exemplary**):
```
RANDOM_ORG_API_KEY=af83r3m2-mv82-z327-12m9238hjqdn
CLIENT_EMAIL=testOwner@gmail.com
PORT=8080
```


## :arrow_up: Room for improvement

- more tests


## :construction_worker: Author

- [@BSski](https://www.github.com/BSski)

</p>
