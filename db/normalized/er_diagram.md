[![](https://mermaid.ink/img/pako:eNq1VNFu0zAU_RXLz1nVNClp8oa0gqLRUbXjBUWKvPiutZbYke2wlawSbzzyBXzGJECIb0j_CCdtWZdlm0DCD45yzvW95_pe3xInggIOMMhjRhaSZBFHZk3Hs_nbU3Rzc3Qkyt1fPJ2NJ-G7CQrQkqhDuz9MuUXrxbh2EaNoetLCcpBK8NhQrw4oSjQ0W6w0kUBRFwWcAt0S63s6n4mr4VojTjJoOT1nUi8pWaGWKQWVSJZrJniLSUVCOmDICEtbWE6UuhKSHsAsA5NdlqNEgolPY6JbwlN2CSpO4eKAWNRpS7X7dqV_v0zh5OXrcUeRtvi_l6ihEkjTmBfZOUgU4epr9av6Wd1uPqHNl83n6rb6Xv2ovkW4dRcsIwuIC5k-of44nL8JT2rhEV6AVohwihbsA6i9u71F2S1XcIibbKyHsncm-kq0TTq1PCPkv6kIT8_Gs_H8rF3RPRwgplqWf9f6T_bNXZQHrfNYtLtO6UzYbCBNzz-SbsSxhTOQ5vVQM4Ma5xHWSzBycX39lMjL-tLXxo4UWsxXPMGBlgVYuMjrN7ybWji4IKkyaE44Dkp8jQPb73nuwHH8gWv3Pb8_svDKoC-8njscuF7fdYbe0PHXFv4ohHFg9xzbcxzPG9rOaDR0fXMAKNNCTrYjspmUTYT3zYFaxvo3KUCQAQ?type=png)](https://mermaid.live/edit#pako:eNq1VNFu0zAU_RXLz1nVNClp8oa0gqLRUbXjBUWKvPiutZbYke2wlawSbzzyBXzGJECIb0j_CCdtWZdlm0DCD45yzvW95_pe3xInggIOMMhjRhaSZBFHZk3Hs_nbU3Rzc3Qkyt1fPJ2NJ-G7CQrQkqhDuz9MuUXrxbh2EaNoetLCcpBK8NhQrw4oSjQ0W6w0kUBRFwWcAt0S63s6n4mr4VojTjJoOT1nUi8pWaGWKQWVSJZrJniLSUVCOmDICEtbWE6UuhKSHsAsA5NdlqNEgolPY6JbwlN2CSpO4eKAWNRpS7X7dqV_v0zh5OXrcUeRtvi_l6ihEkjTmBfZOUgU4epr9av6Wd1uPqHNl83n6rb6Xv2ovkW4dRcsIwuIC5k-of44nL8JT2rhEV6AVohwihbsA6i9u71F2S1XcIibbKyHsncm-kq0TTq1PCPkv6kIT8_Gs_H8rF3RPRwgplqWf9f6T_bNXZQHrfNYtLtO6UzYbCBNzz-SbsSxhTOQ5vVQM4Ma5xHWSzBycX39lMjL-tLXxo4UWsxXPMGBlgVYuMjrN7ybWji4IKkyaE44Dkp8jQPb73nuwHH8gWv3Pb8_svDKoC-8njscuF7fdYbe0PHXFv4ohHFg9xzbcxzPG9rOaDR0fXMAKNNCTrYjspmUTYT3zYFaxvo3KUCQAQ)

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
    
