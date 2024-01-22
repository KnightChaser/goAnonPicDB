# goAnonPicDB

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![MySQL](https://img.shields.io/badge/mysql-%2300f.svg?style=for-the-badge&logo=mysql&logoColor=white)
![Bootstrap](https://img.shields.io/badge/bootstrap-%238511FA.svg?style=for-the-badge&logo=bootstrap&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?style=for-the-badge&logo=ubuntu&logoColor=white)

**goAnonPicDB** is a simple web service that allows users to anonymously upload images. It is built with Go and uses `GORM` as the **ORM (Object-Relational Mapping)** library. The images are stored in a MySQL database. It might be a good example for learning how to organize applications with `docker` and `docker-compose`.

> **Preview**
<br>User interface (web)
<img src="./readme_pictures/webpage.png">
<br>Database managing uploaded pictures via filepath (real image data are uploaded to `/static`)
<img src="./readme_pictures/db.png">

## Feature

* Anonymous image uploading via HTTP communication
* Display of the 6 most recent images (number may be adjusted by user. Go to `main.go`)
* Count of total uploaded images
* Docker support for easy deployment (via `docker-compose.yaml`)

## Usage

1. **Clone** the repository:

    ```bash
    git clone https://github.com/KnightChaser/goAnonPicDB.git
    cd goAnonPicDB
    ```

2. (optional) Modify `docker-compose.yaml` for preconfigured environmental variables.

    ```yaml
    version: '3'
    services:
    mysql:
        image: mysql:8.0
        environment:
        MYSQL_ROOT_PASSWORD: root   # root password
        MYSQL_DATABASE: images      # db name for images (should be correspondned to main.go init() config)
        MYSQL_USER: knightchaser    # non-root username
        MYSQL_PASSWORD: pass12##    # non-root password
        networks:
        internal_network:
            ipv4_address: 192.168.1.77

    backend:
        build:
        context: .
        environment:
        MYSQL_USERNAME: knightchaser # non-root username
        MYSQL_PASSWORD: pass12##     # non-root password
        networks:
        internal_network:
            ipv4_address: 192.168.1.100
        ports:
        - "8080:8080"

    networks:
    internal_network:
        ipam:
        config:
            - subnet: 192.168.1.0/24
            gateway: 192.168.1.1
    ```

3. Build and run the Docker containers:

    ```bash
    docker-compose up -d
    ```

4. Access the application at [http://localhost:8080](http://localhost:8080) in your web browser. And you'll see the web implementation of `goAnonPicDB`!


## To-Do
- [x] Dockerize applications
- [ ] Create `.env` file and manage cross-service environmental variables and constants (Or do the same thing with other methods)

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.