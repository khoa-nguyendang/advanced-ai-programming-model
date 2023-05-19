use app;

drop procedure if exists sp_device_subscribeack;
delimiter //
create procedure sp_device_subscribeack(
    in deviceUUID varchar(250), 
    in companyCode varchar(250), 
    in subscribeTopic varchar(250), 
    in subscribeState int, 
    in lastModified bigint)
    begin
        if exists(select id from mqtt_subscriber where device_uuid like deviceUUID and company_code like companyCode and topic like topic) then
			update mqtt_subscriber set subscribe_state = subscribeState
            where evice_uuid like deviceUUID
            and company_code like companyCode
            and topic like topic;
        
        else
		    insert into mqtt_subscriber(device_uuid, company_code, topic, subscribe_state, last_modified)
            values(deviceUUID, companyCode, subscribeTopic, subscribeState, lastModified);
            set result = LastInsertedId();
        end if;
    end //
delimiter ;


drop procedure if exists sp_device_GetSubscribers;
delimiter //
create procedure sp_device_GetSubscribers(
    in searchTerm varchar(250), 
    in take int, 
    in skip int)
    begin
        select * 
        from mqtt_subscriber
        where device_uuid like searchTerm or company_code like searchTerm
        limit take, skip;
    end //
delimiter ;
