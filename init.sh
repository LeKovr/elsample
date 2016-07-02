# consup image init file
# run at container start

CONF=/etc/supervisor/conf.d/app.conf
[ -e $CONF ] || ln -s /home/app/supervisor.conf $CONF

echo "Server http://$NODENAME ready"
