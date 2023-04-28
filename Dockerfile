# Official RabbitMQ image
FROM rabbitmq:3.11-management

# Open ports 5672 (AMQP) and 15672 (RabbitMQ Management)
EXPOSE 5672 15672

# Adding a user and setting up access rights
RUN rabbitmqctl add_user USER PASSWORD \
 && rabbitmqctl set_user_tags USER administrator \
 && rabbitmqctl delete_user guest \
 && rabbitmqctl add_vhost customers \
 && rabbitmqctl set_permissions -p customers USER ".*" ".*" ".*"

# Running RabbitMQ
CMD ["rabbitmq-server"]
