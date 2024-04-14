```mermaid
erDiagram
    PERSON ||--o{ PERSON_PREMIUM : has
    PERSON_PREMIUM {
        int4 id PK
        int4 person_id FK
        date date_stared 
        date date_ended
    }
    PERSON {
        int4 id PK
        text name
        date birthday 
        text description
        text location
        text email
        text password
        timestamp created_at
        int4 likes_left
        genders gender
    }
    PERSON ||--o{ PERSON_IMAGE : has
    PERSON_IMAGE {
        int4 id PK
        int4 person_id FK
        int4 cell_number "Номер ячейки"
        text image_url
    }
    PERSON ||--o{ DISLIKE : "gets and gives"
    DISLIKE {
        int4 person_one_id PK, FK
        int4 person_two_id PK, FK
    }
    PERSON ||--o{ LIKE : "gets and gives"
    LIKE {
        int4 person_one_id PK, FK
        int4 person_two_id PK, FK
    }
    INTEREST ||--o{ PERSON_INTEREST : is
    INTEREST {
        int4 id PK
        text name
    }
    PERSON ||--o{ PERSON_INTEREST : has
    PERSON_INTEREST {
        int4 person_id PK, FK
        int4 interest_id PK, FK
    }
