[![Open in Visual Studio Code](https://classroom.github.com/assets/open-in-vscode-2e0aaae1b6195c2367325f4f02e2d04e9abb55f0b24a779b69b11b9e10269abc.svg)](https://classroom.github.com/online_ide?assignment_repo_id=21658079&assignment_repo_type=AssignmentRepo)
[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/GMrD03Jz)
# Graded Challenge 1 - P3

Graded Challenge ini dibuat guna mengevaluasi pembelajaran pada Hacktiv8 Program Fulltime Golang khususnya pada pembelajaran mongoDB dan implementasi microservice golang

## Assignment Objectives
Graded Challenge 1 ini dibuat guna mengevaluasi pemahaman mongoDB dan Microservice sebagai berikut:

- Mampu memahami konsep Microservice
- Mampu memahami konsep Mongo DB
- Mampu mengimplementasikan Mongo DB ke REST APi Golang

## Assignment Directions: Restful API
Buatlah sebuah REST API dengan menggunakan Echo Golang dan implementasikan database MongoDB sesuai dengan kriteria yang telah ditentukan
1. Buat service <strong>Payment</strong> yang memiliki 1 endpoint sebagai berikut:
	- /payments (POST) - Create a new payment
2. 
3. Buat service <strong>Shopping</strong> yang terintegrasi dengan service payment yang memiliki endpoint sebagai berikut:
	- /transactions (POST) - Create a new transactions
	- /transactions (GET) - Retrieve all transactions
	- /transactions/{id} (GET) - Retrieve a specific transactions by ID
	- /transactions/{id} (PUT) - Update a transaction informations
	- /transactions/{id} (DELETE) - Delete a transaction informations
    - /products (POST) - Create a new product
    - /products (GET) - Retrieve all products
    - /products/{id} (GET) - Retrieve a specific product by ID
    - /products/{id} (PUT) - Update a product informations
    - /products/{id} (DELETE) - Delete a product informations
> Pada Endpoint /transactions (POST) Sambungkan dengan service payment untuk melakukan pembayaran, Jika pembayaran berhasil maka transaksi berhasil, jika gagal maka transaksi juga gagal

Gunakan function untuk handle error seperti yang dilakukan pada phase 2, Kalian harus menggunakan mongodb untuk menyimpan dan melakukan transaksi penyewaan lapangan.

#### Sebagai tambahan dari requirement yang sudah diberikan sebelumnya, Student juga diharapkan untuk memahami dan menerapkan konsep-konsep berikut:
- Cloud Deployment using GCP
Student diharapkan untuk mengimplementasikan Cloud Deployment menggunakan Google Cloud Platform (GCP).
Pastikan aplikasi Anda dapat diakses secara publik setelah deployment.
Sediakan dokumentasi sederhana mengenai langkah-langkah deployment yang Anda lakukan.
- Job Scheduling
Implementasikan konsep job scheduling untuk beberapa proses yang memerlukannya, seperti proses pembaharuan data atau pembersihan data yang tidak diperlukan.
- Unit Test
Buat unit test untuk memastikan bahwa setiap fungsi atau method dalam aplikasi Anda bekerja dengan semestinya.
- Docker
Kontainerisasi aplikasi Anda menggunakan Docker.
Pastikan Anda menyediakan Dockerfile dan dokumentasi singkat tentang bagaimana menjalankan aplikasi Anda menggunakan Docker.

## Database Schema
> Cobalah buat Database Schema dari analisa kebutuhan endpoint diatas, untuk nama table sebagai berikut :
>  - Products
>  - Transactions
>  - Payments
## Expected Results

- Aplikasi terdiri dari 2 :
- Service Shopping berjalan pada http://localhost:8080 dengan endpoint berikut :
  - /transactions (POST) - Create a new transactions
  - /transactions (GET) - Retrieve all transactions
  - /transactions/{id} (GET) - Retrieve a specific transactions by ID
  - /transactions/{id} (PUT) - Update a transaction informations
  - /transactions/{id} (DELETE) - Delete a transaction informations
  - /products (POST) - Create a new product
  - /products (GET) - Retrieve all products
  - /products/{id} (GET) - Retrieve a specific product by ID
  - /products/{id} (PUT) - Update a product informations
  - /products/{id} (DELETE) - Delete a product informations
  
- Service Payment berjalan pada http://localhost:8081
  - /payments (POST) - Create a new payment


### RESTRICTION:
- The id field should be the primary key for the table.
- The email field should have a unique index to prevent duplicate email addresses.


---------- 

###  Assignment Submission

Push Assigment yang telah Anda buat ke akun Github Classroom Anda masing-masing.

----------

## Assignment Rubrics

Aspect : 
|Criteria|Meet Expectations|Points|
|---|---|---|
|Problem Solving|API Endpoints are implemented and working correctly |75 pts |
|Database Design |MongoDB database meets the required specifications |10 pts|
||Database queries are efficient and appropriately indexed |5 pts|
|Readability|Code is well-documented and easy to read |5 pts|
||Code includes appropriate comments and documentation |5 pts|


Notes:
- Don't rush through the problem or try to solve it all at once.
- Don't copy and paste code from external sources without fully understanding how it works.
- Don't hardcode values or rely on assumptions that may not hold true in all 
cases.
- Don't forget to handle error cases or edge cases, such as invalid input or unexpected behavior.
- Don't hesitate to refactor your code or make improvements based on feedback or new insights.



Total Points : 100

Notes Deadline : W2D1 - 9.00AM

Keterlambatan pengumpulan tugas mengakibatkan skor GC 1 menjadi 0.
