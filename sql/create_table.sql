CREATE TABLE all_messages ( 
    id SERIAL, 
    timestamp TIMESTAMP, 
    topic VARCHAR(255), 
    title VARCHAR(255),     
    comment TEXT, 
    sendtokafka BOOLEAN
);

CREATE USER m_user WITH PASSWORD 'sdvFGGSV456asd';
GRANT ALL ON all_messages TO m_user;
GRANT ALL ON all_messages_id_seq TO m_user;