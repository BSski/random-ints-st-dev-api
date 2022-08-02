<p align="center">
  <h1 align="center">Random Integers' Standard Deviation API</h2>
</p>


<div markdown="1" align="center">    

![Demo Screenshot](https://i.imgur.com/OgMTzZX.png)

</div>


## Table of contents
* [Project description](#scroll-project-description)
* [Technologies used](#hammer-technologies-used)
* [Deployment](#hammer_and_wrench-deployment)
* [Environment variables](#closed_lock_with_key-environment-variables)
* [Room for improvement](#arrow_up-room-for-improvement)
* [Author](#construction_worker-author)


## :scroll: Project description
This is a recruitment task for a Junior Software Developer position in Go.

The project is a REST service supporting the following GET operation:

```/random/mean?requests={r}&length={l}```

which performs `{r}` concurrent requests to random.org API asking for `{l}` number of random integers.

Additionally, the application calculates standard deviation of each drawn integers set and of the sum of all sets.

The application is also dockerized.


## :hammer: Technologies used
- Go 1.18
- go-chi/chi/v5 5.0.7
- google/go-cmp 0.5.8
- joho/godotenv 1.4.0
- montanaflynn/stats 0.6.6
- cosmtrek/air 1.40.4
- Docker


## :hammer_and_wrench: Deployment

1. Create an `.env` file basing on `.env_sample_file` from the repository. Set `PORT` to 8080.

2. Run `docker run --env-file .env -p 8080:8080 bsski/random-ints-st-dev-api:latest` in the `.env` file directory.

3. Access `localhost:8080`. 


## :closed_lock_with_key: Environment variables

To run this project, you have to set up the following environment variables in the `.env` file (**the values below are exemplary**):
```
RANDOM_ORG_API_KEY=af83r3m2-mv82-z327-12m9238hjqdn
PORT=8080
```


## :arrow_up: Room for improvement

- more and better tests,
- CORS middleware.


## :construction_worker: Author

- [@BSski](https://www.github.com/BSski)
