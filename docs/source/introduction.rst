What is Waitress?
=================

Waitress is an innovative platform designed to enhance the dining experience by integrating food delivery and restaurant reservation systems. 
It enables customers to effortlessly reserve their preferred table or seat at their favorite restaurant, 
pre-order their meal via a web browser or mobile phone, and specify their arrival time. 
This seamless process ensures that guests are seated and served within 10 to 15 minutes of their arrival, 
merging convenience with a quality dining experience.

Purpose
-------

The goal of Waitress is to bridge the gap between traditional dining and modern technological convenience, 
offering a unique solution for both restaurants and diners who would like to dine at the restaurant. 

Key Features
------------

* **Table Reservation**: Book your favorite table or seat in advance.
* **Pre-Order Meals**: Select and order your meal ahead of your visit.
* **Scheduled Dining**: Choose your arrival time to minimize waiting.
* **Swift Service**: Be seated and served within 10 - 15 minutes of arrival.

Getting Started
---------------

Flask Configuration Setup
~~~~~~~~~~~~~~~~~~~~~~~~~~

By now, you should have the repository cloned down and ready to go. Please ensure you message `Rahmin Shoukoohi <https://github.com/1ShoukR/>`_ 
for a copy of a ``local.cfg`` file. This file is needed for running the file locally. NEVER commit this file to Git (it is ignored byt default in ``.gitignore``)

Database Setup
~~~~~~~~~~~~~~~

Create a a new database using MySQL Workbench (recommended name is *waitress*)
You may need to adjust values in ``local.cfg`` to match the MySQL settings on your machine.
MySQL workbench does not setup MySQL consistently, especially across different operating systems, so this is hard to predict.
Once you figure out your credentials and which port the MySQL service is running on, update ``SQLALCHEMY_DATABASE_URI`` in ``local.cfg``::

    # Format
    # SQLALCHEMY_DATABASE_URI = 'mysql+pymysql://{username}:{password}@localhost:{port}/{dbname}?charset=utf8mb4'

    # Actual Example
    SQLALCHEMY_DATABASE_URI = 'mysql+pymysql://root:root@localhost:3306/sym?charset=utf8mb4'

There currently is no mock data to implement and seed into the database, however, this will be changed as that is next on the priority list