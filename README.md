# uber-replica
A uber replica (Customer and Driver app) with Golang (Gin) and VueJS and Bizops app with React & TailwindCSS

## How to run:
for server osrm and db: `docker-compose up --build`
(Container osrm need to generate data at `backend/data`, try on this [document](https://gist.github.com/AlexandraKapp/e0eee2beacc93e765113aff43ec77789))

for backend: `./backend/go run main.go` (on packaging to container)

for customer and driver frontend: `./frontend/npm run serve` (on packaging to container)

for bizops app: `./bizops/npm start` (on packaging to container)

## Customer app
#### Customer main
![alt text](./img/customer-main.png)
#### Customer input address
![alt text](./img/customer-address.png)
#### Customer map and pricing
![alt text](./img/customer-map.png)
#### Customer waiting for driver accept
![alt text](./img/customer-driver.png)

## Driver app
#### Driver main (On or Off accepting)
![alt text](./img/driver-main.png)
#### Driver come to customer
![alt text](./img/driver-to-cus.png)
#### Driver drive to destination
![alt text](./img/driver-to-des.png)
#### Summarise trip
![alt text](./img/sumarize.png)

## Bizops
#### Dashboard
![bizops image](./img/bizops.png)
#### List trips
![list trip](./img/bizops-trip.png)