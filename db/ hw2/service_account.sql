CREATE USER jimder_service_account WITH PASSWORD 'iamoutoftouch888';

CREATE ROLE jimder_service_account_role WITH LOGIN;

GRANT USAGE ON SCHEMA public TO jimder_service_account_role;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO jimder_service_account_role;
GRANT INSERT ON person_interest, person_payment, person, message, person_claim, "like", person_image TO jimder_service_account_role;
GRANT DELETE ON person_interest, person, person_image TO jimder_service_account_role;
GRANT UPDATE ON person TO jimder_service_account_role;
