CREATE USER jimder_service_account WITH PASSWORD 'iamoutoftouch888';

CREATE ROLE jimder_service_account_role WITH LOGIN;

GRANT USAGE ON SCHEMA public TO jimder_service_account_role;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO jimder_service_account_role;
GRANT INSERT ON person, interest, person_interest, "like", dislike, 
person_image, person_premium, "message", person_claim, 
claim, person_payment TO jimder_service_account_role;
GRANT DELETE ON 