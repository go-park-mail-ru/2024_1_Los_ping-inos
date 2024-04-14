# Описание схемы

## Функциональные зависимости
* **Таблица person:** \
{id} -> {name, birthday, description, location, email, password, created_at, likes_left, gender}
{email} -> {name, birthday, descritpion, location, password, created_at, likes_left, gender}

* **Таблица like** \
В данной таблице нет неключевых атрибутов

* **Таблица dislike** \
В данной таблице нет неключевых атрибутов


* **Таблица person_image** \
{person_id, cell_number} -> {image_url}


* **Таблица interest** \
{id} -> {name}

* **Таблица person_interest** \
В данной таблице нет неключевых атрибутов


* **Таблица person_premium** \
{id} -> {person_id, date_started, date_ended}


## Описание 
* **Таблица person:** \
Таблица пользователей, в которой хранятся пароль и логин пользователя, а так же информация о нем. \
  - name - имя пользователя
  - birthday - день рождения пользователя
  - description - описание профиля пользователя, которое отображается в его карточке
  - location - местоположение пользователя
  - email - email пользователя
  - password - пароль пользователя
  - created_at - время создания аккаунта
  - likes_left - количество лайков, оставшееся у пользователя
  - gender - пол пользователя \
  
Данная таблица соответствует всем нужным НФ потому что каждый котреж содержит только одно значение для каждого из атрибутов, нет зависимостей от части ключа, нет функциональных зависимостей между неключевыми атрибутами а также все детерминанты являются потенциальными ключами

* **Таблица like** \
Связующая таблица между двумя пользователями. Она нужна для того чтобы определить произошел ли метч при лайке или нет. Если первый пользователь ставит лайк второму, а в таблице \
есть запись с лайком от второго первому, происходит метч.
  - person_one_id - айди первого пользователя
  - person_two_id - айди второго пользователя \
 
Данная таблица соответствует всем нужным НФ потому что она содержит только ключевые атрибуты


* **Таблица dislike** \
Связующая таблица между двумя пользователями. Она нужна для того чтобы при дизлайке пользователь не попадался снова в ленте.
  - person_one_id - айди первого пользователя
  - person_two_id - айди второго пользователя
 
Данная таблица соответствует всем нужным НФ потому что она содержит только ключевые атрибуты

* **Таблица person_image** \
Таблица, в которой хранится информация о всех фото пользователя. Каждый пользователь может загрузить 5 фото в 5 ячеек. \
Следовательно у каждой картинки есть номер ячейки в которую она загружена. Как первичный ключ она хранит айди пользователя и ячейку картинки.
  - person_id - айди пользователя
  - cell_number - номер ячейки
  - image_url - урл картинки
 
Данная таблица соответствует всем нужным НФ потому что каждый котреж содержит только одно значение для каждого из атрибутов, нет зависимостей от части ключа, нет функциональных зависимостей между неключевыми атрибутами а также все детерминанты являются потенциальными ключами 

* **Таблица interest** \
Таблица, хранящая информацию о всех возможных интересах, которые могут выбрать пользователи.
  - id - id интереса
  - name - название интереса
 
Данная таблица соответствует всем нужным НФ потому что каждый котреж содержит только одно значение для каждого из атрибутов, нет зависимостей от части ключа, нет функциональных зависимостей между неключевыми атрибутами а также все детерминанты являются потенциальными ключами 

* **Таблица person_interest** \
Связующая таблица между пользователем и его интересами.
  - person_id - айди пользователя
  - interest_id - айди интереса
 
Данная таблица соответствует всем нужным НФ потому что она содержит только ключевые атрибуты

* **Таблица person_premium** \
Таблица, в которой хранится информация о премиум подписке пользователя. Запись в ней создается тогда, когда пользователь покупает подписку.
- id - айди подписки
- person_id - айди пользователя
- date_started - дата покупки подписки
- date_ended - дата, до которой подписка действует

Данная таблица соответствует всем нужным НФ потому что каждый котреж содержит только одно значение для каждого из атрибутов, нет зависимостей от части ключа, нет функциональных зависимостей между неключевыми атрибутами а также все детерминанты являются потенциальными ключами 