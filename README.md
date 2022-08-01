<p align="center">
  <h2 align="center">DjangoImagesAPI [work in progress]</h2>
</p>


<div markdown="1" align="center">


[![Build Status](https://bsski.semaphoreci.com/badges/DjangoImagesAPI/branches/main.svg?style=shields&key=0adba4be-b67c-416d-802b-90f43f19ff0f)](https://bsski.semaphoreci.com/projects/DjangoImagesAPI)
[![Demo Uptime](https://img.shields.io/uptimerobot/ratio/7/m792086829-c54e14cd8cfdacdfdfa92920)](https://django-images-api-bsski.herokuapp.com/)
[![Heroku](https://pyheroku-badge.herokuapp.com/?app=django-images-api-bsski&style=flat)](https://django-images-api-bsski.herokuapp.com/)
[![Maintainability](https://api.codeclimate.com/v1/badges/a84055e6e49ddb02653e/maintainability)](https://codeclimate.com/github/BSski/DjangoImagesAPI/maintainability)
[![CodeFactor](https://www.codefactor.io/repository/github/bsski/djangoimagesapi/badge)](https://www.codefactor.io/repository/github/bsski/djangoimagesapi)
[![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/psf/black)
</div>


<h3 align="center">
  :fireworks: Live demo :fireworks:
</h3>

<p align="center">
  https://django-images-api-bsski.herokuapp.com/
</p>

<p align="center">
Read-only admin and test users credentials are available upon request.
</p>

<p align="center">
Login page URL:<br>
https://django-images-api-bsski.herokuapp.com/auth/login
</p>



## Table of contents
* [Project description](#scroll-project-description)
* [Technologies used](#hammer-technologies-used)
* [Deployment](#hammer_and_wrench-deployment)
* [Environment variables](#closed_lock_with_key-environment-variables)
* [Main features](#rocket-main-features)
* [Room for improvement](#arrow_up-room-for-improvement)
* [Author](#construction_worker-author)


## :scroll: Project description
This is a recruitment task for a Junior Python Developer position.

The project:
- is an Images API,
- where users can list their images,
- <a href="https://github.com/BSski/DjangoImagesAPI/blob/development/website/fixtures.json">with 3 main built-in user tiers (Basic, Premium, Enterprise) with varying permissions </a> to:
    - fetch a link to a permitted thumbnail size,
    - fetch an expiring link to a permitted thumbnail size (300-30000 sec),
    - fetch a link to the original image,
- where user tiers with arbitrary permissions can be created via admin panel,
- and if a user tier settings change, all its users are updated, and so are their images,
- and if a user's user tier changes, his images get updated,
- and if an image's owner changed, its thumbnail links are updated,
- user inputted data is validated (includes regex!),
- project's design and architecture was created with performance in mind:
    - only original user's images are stored permanently: thumbnails are stored in a S3 bucket which deletes files that were not accessed for 7 days,
    - thumbnail is generated only when requested, by AWS Lambda function,
    - JWT token is used for authentication to minimize number of connection to the database,
- many throttling mechanisms were added,
- pagination was added,
- HyperlinkedRelatedField for user's info in images list view was added,
- the code has docstrings,
- and it's easily deployable via docker, but you need your own AWS architecture (Localstack will be added soon),
- there is a CI/CD pipeline which would run tests if I wrote them, it runs black linter check, dockerizes the project and deploys it to Heroku from the Docker container,
- there are authentication and authorization mechanisms: user's JWT token is saved as cookie, but it can also be passed in requests as a bearer token.


## :hammer: Technologies used
- Python 3.7 & 3.8
- Django 3.2
- Django REST Framework
- PostgreSQL 14.2
- AWS S3
- AWS Lambda
- Gunicorn
- Docker
- SemaphoreCI
- Heroku
- dj-rest-auth with JWT Token


## :hammer_and_wrench: Deployment


A) through Docker:

1. Create an `.env` file basing on `.env_sample_file` from the repository. Set `PORT` to 8020.

2. Run `docker run --env-file .env -p 8020:8020 bsski/images-api:latest` in the `.env` file directory.

3. Access `localhost:8020`. 


B) without Docker:

1. Download the repository.

2. Create a virtual environment.

3. Run `pip install -r requirements.txt` (or `requirements-windows.txt`) in the directory of `requirements.txt`.

4. Create an `.env` file basing on `.env_sample_file` from the repository in the directory of `.env_sample_file`. Set `PORT` to 8000.

5. Run `python manage.py runserver` in the directory of `manage.py`.

6. Access `127.0.0.1:8000`.


The admin panel can be found under your chosen URL or, if you didn't set it in the .env file, the default `/hidden_admin_url`.

Security through obscurity is not enough of course, but I find it a nice complementary solution.


## :closed_lock_with_key: Environment variables

To run this project, you have to set up the following environment variables in the `.env` file (**the values below are exemplary**):
```
SECRET_KEY=TestSecretKey
AWS_S3_REGION_NAME=eu-central-1
AWS_S3_ADDRESSING_STYLE=virtal
AWS_S3_SIGNATURE_VERSION=s3v4
AWS_ACCESS_KEY_ID=MFAA21AFHG9AFKA2AEY
AWS_SECRET_ACCESS_KEY=13Fjm32eam23k23Qf432fDFda9saf
AWS_STORAGE_BUCKET_NAME=main-bucket
AWS_THUMBNAILS_STORAGE_BUCKET_NAME=thumbnails-bucket
DEFAULT_FILE_STORAGE=storages.backends.s3boto3.S3Boto3Storage
PORT=8020
ADMIN_LOGIN_URL=test_admin_login_url
DJANGO_SUPERUSER_PASSWORD=TestAdminPassword
DJANGO_SUPERUSER_USERNAME=TestAdminName
DJANGO_SUPERUSER_EMAIL=TestAdmin@email.com
LOCALSTACK_ENDPOINT_URL=http://localhost:4566
```


## :rocket: Main features

The recruitment task demanded such features and all are provided in a required format:
- sending a POST request to `/auth/login/` lets you log in,
- sending a POST request to `/auth/logout/` lets you log out,
- accessing `/images/images/` when logged in lists current user's images,
- accessing 'images/images/' when logged in with `?limit=3` at the end of it sets a pagination for the current page,
- accessing thumbnail link from image view with `?time_exp=500` at the end of it generates a link to it expiring in 500 seconds (if logged in user is permitted to do that; range: 300-30000).


Furthermore, the website is deployed on Heroku from a Docker image using a CI/CD SemaphoreCI pipeline:

![CI/CD screenshot](https://i.imgur.com/3kWt2aT.png)


## :arrow_up: Room for improvement

Due to the time limitation, there are very important things left to add:
- [in progress] Localstack to test and run app locally,
- [in progress] tests,
- typehints,
- global AWS clients instead of temporary ones,
- the code would also benefit from some cleaning,
- logging.


## :construction_worker: Author

- [@BSski](https://www.github.com/BSski)
