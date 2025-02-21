--Create DB จากตรงนี้เป็นต้นไป
CREATE TABLE IF NOT EXISTS customers(
    Cust_ID SERIAL  PRIMARY KEY,
    Cust_Fname VARCHAR(50) NOT NULL,
    Cust_Lname VARCHAR(50) NOT NULL,
    Cust_Email VARCHAR(50) DEFAULT NULL,
    Cust_PhoneNumber VARCHAR(10) NOT NULL,
    Cust_Username VARCHAR(50) NOT NULL,
    Cust_Password VARCHAR(50) NOT NULL
);
CREATE TABLE IF NOT EXISTS cust_address(
    PostlalCode VARCHAR(5) PRIMARY KEY NOT NULL,
    Cust_ID SERIAL NOT NULL,
    House_No VARCHAR(10),
    Province VARCHAR(25),
    Village_NO VARCHAR(5),
    Soi VARCHAR(25),
    District VARCHAR(25),
    Sub_District VARCHAR(25),
    Road VARCHAR(25),
    FOREIGN KEY(Cust_ID) REFERENCES customers(Cust_ID) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS shops(
    Shop_ID SERIAL PRIMARY KEY,
    Shop_Name VARCHAR(50) NOT NULL,
    Shop_Des text,
    Shop_Image VARCHAR(255),
    Shop_Address VARCHAR(255),
    Shop_PhoneNumber VARCHAR(10),
    Shop_Email VARCHAR(50)
);
CREATE TABLE IF NOT EXISTS room_type(
    Room_ID SERIAL PRIMARY KEY,
    Room_Name VARCHAR(50)
);
CREATE TABLE IF NOT EXISTS fur_type(
    Fur_ID SERIAL PRIMARY KEY,
    Fur_Name VARCHAR(50)
);
CREATE TABLE IF NOT EXISTS brand(
    Brand_ID SERIAL PRIMARY KEY,
    Brand_Name VARCHAR(50) DEFAULT 'No name',
    Brand_Image VARCHAR(255)
);
CREATE TABLE IF NOT EXISTS products(
    Prod_ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    Prod_Name VARCHAR(100) NOT NULL,
    Prod_Price NUMERIC(10, 2) NOT NULL CHECK (Prod_Price >= 0),
    Prod_Details text,
    Prod_Image VARCHAR(255),
    Stock NUMERIC NOT NULL CHECK (Stock >= 0),
    Sales_Amount NUMERIC NOT NULL CHECK (Sales_Amount >= 0),
    Status_Rec BOOLEAN DEFAULT FALSE,
    Created_At TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    Updated_At TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    Room_ID SERIAL NOT NULL,
    fur_ID SERIAL NOT NULL,
    Brand_ID SERIAL Not NULL,
    Shop_ID SERIAL NOT NULL,
    FOREIGN KEY(Brand_ID) REFERENCES brand(Brand_ID) ON DELETE SET DEFAULT,
    FOREIGN KEY(Room_ID) REFERENCES room_type(Room_ID) ON DELETE SET DEFAULT,
    FOREIGN KEY(Fur_ID) REFERENCES fur_type(Fur_ID) ON DELETE SET DEFAULT,
    FOREIGN KEY(Shop_ID) REFERENCES shops(Shop_ID) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS choosing(
    Cust_ID SERIAL NOT NULL,
    Prod_ID UUID DEFAULT uuid_generate_v4() NOT NULL,
     PRIMARY KEY (Cust_ID, Prod_ID),
    FOREIGN KEY(Cust_ID) REFERENCES customers(Cust_ID) ON DELETE CASCADE,
    FOREIGN KEY(Prod_ID) REFERENCES products(Prod_ID) ON DELETE CASCADE
);

INSERT INTO customers ( Cust_Fname, Cust_Lname, Cust_Email, Cust_PhoneNumber, Cust_Username, Cust_Password ) 
VALUES ( 'John', 'Doe', 'john.doe@example.com', '0812345678', 'johndoe' ,'password123'),
('Alice', 'Smith', 'alice.smith@example.com', '0898765432', 'alicesmith', 'alicepass'), 
('Bob', 'Brown', 'bob.brown@example.com', '0854321987', 'bobbrown', 'bobpass123'),
 ('Charlie', 'Johnson', 'charlie.j@example.com', '0811122233', 'charliejohn', 'charliepw'), 
('Diana', 'Lee', 'diana.lee@example.com', '0822233445', 'dianalee', 'dianapw456');
INSERT INTO cust_address ( PostlalCode, Cust_ID, House_No, Province, Village_NO, Soi, District, Sub_District, Road ) 
VALUES ('10110', 1, '123', 'Bangkok', '4', 'Soi 1', 'Pathum Wan', 'Lumphini', 'Wireless Road'), 
('10200', 2, '456', 'Bangkok', '3', 'Soi 2', 'Dusit', 'Wachira Phayaban', 'Samsen Road'), 
('10330', 3, '789', 'Bangkok', '2', 'Soi 3', 'Ratchathewi', 'Thanon Phaya Thai', 'Phaya Thai Road'),
 ('10400', 4, '159', 'Bangkok', '1', 'Soi 4', 'Chatuchak', 'Chatuchak', 'Kamphaeng Phet 2 Road'),
 ('10500', 5, '753', 'Bangkok', '5', 'Soi 5', 'Bang Rak', 'Si Lom', 'Silom Road');
--ข้อมูลร้านค้าแต่ละร้าน 
INSERT INTO shops ( Shop_ID, Shop_Name, Shop_Des, Shop_Image, Shop_Address, Shop_PhoneNumber, Shop_Email ) 
VALUES (1, 'White & Wood', 'ร้านเฟอร์นิเจอร์มินิมอลที่เน้นลายไม้และสีขาว', 'https://cdn.pixabay.com/photo/2019/06/02/07/23/coffee-shop-4245788_1280.jpg', '123 Main St, Bangkok', '0812345678', 'contact@whiteandwood.com'), 
(2, 'Home Heaven', 'ร้านเฟอร์นิเจอร์และสินค้าของร้านจะทำให้บ้านของลูกค้ากลายเป็นที่พักผ่อนที่เหมือนอยู่บนสรวงสวรรค์', 'https://cdn.pixabay.com/photo/2021/12/23/03/58/da-guojing-6888603_1280.jpg', '456 Home St, Bangkok', '0823456789', 'info@homeheaven.com'), 
(3, 'The Living Space', 'ร้านเฟอร์นิเจอร์ที่เน้นการสร้างบรรยากาศที่อบอุ่นและสวยงามภายในบ้าน', 'https://cdn.pixabay.com/photo/2016/11/19/17/25/furniture-1840463_1280.jpg', '789 Living St, Bangkok', '0834567890', 'support@thelivingspace.com'),
 (4, 'Tranquil Home', 'ร้านเฟอร์นิเจอร์ที่ต้องการสร้างบรรยากาศผ่อนคลายและความเป็นมินิมอลให้กับลูกค้า', 'https://cdn.pixabay.com/photo/2016/11/21/12/59/couch-1845270_1280.jpg', '101 Home St, Bangkok', '0845678901', 'hello@tranquilhome.com'), 
(5, 'EchoHome', 'ร้านเฟอร์นิเจอร์ที่ออกแบบมาเพื่อให้บ้านของคุณมีความทันสมัยและสะดวกสบายที่สุด ด้วยสไตล์มินิมอลที่เน้นการใช้รูปทรงเรียบง่าย สีโทนอบอุ่น และวัสดุธรรมชาติ', 'https://cdn.pixabay.com/photo/2017/09/09/18/25/living-room-2732939_1280.jpg', '202 Pet St, Bangkok', '0856789012', 'contact@echohome.com');
Insert into room_type(Room_Name) values ('Livingroom'), ('Bedroom'), ('Bathroom'), ('Kitchen'), ('Others');
Insert into fur_type(Fur_Name) values ('Sofa'), ('Bed'), ('Bath'), ('Chair'), ('Table'), ('Others');

Insert into brand (Brand_ID, Brand_Name, Brand_Image ) 
VALUES (1, 'IKEE', 'https://cdn.pixabay.com/photo/2023/09/15/12/43/living-room-8254772_1280.jpg'),
(2, 'The Lock', 'https://cdn.pixabay.com/photo/2022/01/20/03/50/icon-6951393_1280.jpg'),
(3, 'Jeff santa','https://cdn.pixabay.com/photo/2022/02/08/01/08/gold-heart-7000551_1280.jpg'),
(4, 'Omai','https://cdn.pixabay.com/photo/2022/03/27/21/29/shark-7096133_1280.png');

-- เพิ่มข้อมูลในตาราง Products โดยแต่ละร้านมี 5 ผลิตภัณฑ์
-- ตัวอย่างข้อมูลสำหรับร้านค้า 5 ร้าน แต่ละร้านมีสินค้า 5 ชิ้น รวมเป็น 25 รายการ
INSERT INTO products (Prod_Name, Prod_Price, Prod_Details, Prod_Image, Stock, Sales_Amount, Room_ID, Fur_ID, Brand_ID, Shop_ID)
VALUES
-- ร้านที่ 1 
('IKEE Livingroom Sofa', 999.99, 'Comfortable sofa for living room.', 'https://cdn.pixabay.com/photo/2017/08/02/01/01/living-room-2569325_1280.jpg', 10, 5, 1, 1, 1, 1),
('IKEE Bedroom Bed', 1299.50, 'Cozy bed for your bedroom.', 'https://cdn.pixabay.com/photo/2016/11/19/13/06/bed-1839183_1280.jpg', 15, 3, 2, 2, 1, 1),
('IKEE Bathroom Bath', 599.00, 'Elegant bath tub.', 'https://cdn.pixabay.com/photo/2022/10/02/14/06/bath-7493560_1280.jpg', 5, 2, 3, 3, 1, 1),
('IKEE Kitchen Chair', 199.75, 'Durable kitchen chair.', 'https://cdn.pixabay.com/photo/2017/06/13/23/14/chair-2400521_1280.jpg', 20, 7, 4, 4, 1, 1),
('IKEE Others Table', 249.99, 'Multipurpose table.', 'https://cdn.pixabay.com/photo/2017/08/01/23/51/apple-2568755_1280.jpg', 12, 4, 5, 5, 1, 1),
-- ร้านที่ 2 (กฤต)
('The Lock Livingroom Sofa', 1050.00, 'Stylish sofa for modern living rooms.', 'https://cdn.pixabay.com/photo/2016/11/23/14/29/living-room-1853203_1280.jpg', 8, 3, 1, 1, 2, 2),
('The Lock Bedroom Bed', 1100.00, 'Comfortable bed with storage.', 'https://cdn.pixabay.com/photo/2018/01/24/15/08/live-3104077_1280.jpg', 10, 5, 2, 2, 2, 2),
('The Lock Bathroom Bath', 650.00, 'Modern bath for small bathrooms.', 'https://cdn.pixabay.com/photo/2018/01/29/07/55/modern-minimalist-bathroom-3115450_1280.jpg', 7, 2, 3, 3, 2, 2),
('The Lock Kitchen Chair', 175.50, 'Compact chair for kitchen use.', 'https://cdn.pixabay.com/photo/2017/03/20/19/57/chairs-2160184_1280.jpg', 18, 3, 4, 4, 2, 2),
('The Lock Others Table', 275.00, 'Table suitable for various uses.', 'https://cdn.pixabay.com/photo/2017/03/28/12/17/chairs-2181994_1280.jpg', 9, 6, 5, 5, 2, 2),  
-- ร้านที่ 3
('Jeff Santa Livingroom Sofa', 1150.00, 'Luxurious sofa from Jeff Santa.', 'https://cdn.pixabay.com/photo/2020/05/24/09/52/sofa-5213406_1280.jpg', 5, 7, 1, 1, 3, 3),
('Jeff Santa Bedroom Bed', 1399.99, 'Elegant bed by Jeff Santa.', 'https://cdn.pixabay.com/photo/2014/08/11/21/40/bedroom-416062_1280.jpg', 12, 5, 2, 2, 3, 3),
('Jeff Santa Bathroom Bath', 799.99, 'Luxurious bath tub.', 'https://cdn.pixabay.com/photo/2010/12/14/14/58/bath-3148_1280.jpg', 6, 1, 3, 3, 3, 3),
('Jeff Santa Kitchen Chair', 225.00, 'Modern kitchen chair.', 'https://cdn.pixabay.com/photo/2013/12/25/16/40/seat-233625_1280.jpg', 11, 8, 4, 4, 3, 3),
('Jeff Santa Others Table', 299.99, 'High-quality table.', 'https://cdn.pixabay.com/photo/2016/11/19/17/25/furniture-1840463_1280.jpg', 8, 2, 5, 5, 3, 3),
-- ร้านที่ 4
('Omai Livingroom Sofa', 999.00, 'Stylish sofa from Omai.', 'https://cdn.pixabay.com/photo/2017/03/19/01/43/living-room-2155376_1280.jpg', 10, 6, 1, 1, 4, 4),
('Omai Bedroom Bed', 1250.00, 'Comfortable bed from Omai.', 'https://cdn.pixabay.com/photo/2018/09/15/09/05/home-3678956_1280.jpg', 9, 3, 2, 2, 4, 4),
('Omai Bathroom Bath', 750.00, 'Modern bath from Omai.', 'https://cdn.pixabay.com/photo/2018/08/09/03/58/home-3593729_1280.jpg', 6, 4, 3, 3, 4, 4),
('Omai Kitchen Chair', 200.00, 'Functional chair for kitchen.', 'https://cdn.pixabay.com/photo/2015/08/19/08/32/beach-895685_1280.jpg', 14, 7, 4, 4, 4, 4),
('Omai Others Table', 275.00, 'Durable table by Omai.', 'https://cdn.pixabay.com/photo/2016/11/18/14/05/brick-wall-1834784_1280.jpg', 12, 1, 5, 5, 4, 4),
-- ร้านที่ 5
('The Lock Livingroom Sofa', 1200.00, 'Comfortable sofa for any living space.', 'https://cdn.pixabay.com/photo/2018/01/20/09/42/sofa-3094153_1280.jpg', 8, 3, 1, 1, 2, 5),
('The Lock Bedroom Bed', 1400.00, 'Spacious bed by The Lock.', 'https://cdn.pixabay.com/photo/2016/04/17/08/09/room-1334323_1280.jpg', 15, 2, 2, 2, 2, 5),
('The Lock Bathroom Bath', 650.00, 'Relaxing bath tub.', 'https://cdn.pixabay.com/photo/2018/07/26/10/36/bathroom-3563272_1280.jpg', 7, 5, 3, 3, 2, 5),
('The Lock Kitchen Chair', 210.50, 'Modern kitchen chair.', 'https://cdn.pixabay.com/photo/2014/02/21/00/09/chair-270980_1280.jpg', 20, 4, 4, 4, 2, 5),
('The Lock Others Table', 250.00, 'Versatile table.', 'https://cdn.pixabay.com/photo/2024/05/13/10/02/dining-table-8758636_1280.jpg', 9, 7, 5, 5, 2, 5);

 --ตาราง ข้อมูลสินค้าในตะกร้า (ของคิว)
CREATE TABLE cart (
    cart_id SERIAL PRIMARY KEY,  
    cust_id INT NOT NULL,
    totalprice NUMERIC(10, 2),
    created_at TIMESTAMP ,  
    updated_at TIMESTAMP ,  
    FOREIGN KEY (cust_ID) REFERENCES customers(cust_id)  
);
CREATE TABLE cart_items (
    prod_image  VARCHAR(255),                         
    prod_name VARCHAR(255),                
    brand_name VARCHAR(255),               
    quantity INT DEFAULT 1,               
    prod_price NUMERIC(10, 2),             
    totalprice NUMERIC(10, 2),
	cart_item_id SERIAL PRIMARY KEY,
	prod_id UUID,
	cart_id INT,
    FOREIGN KEY (prod_id) REFERENCES products(prod_id),
	FOREIGN KEY (cart_id) REFERENCES cart(cart_id) 
);


	   CREATE OR REPLACE FUNCTION update_insert() 
RETURNS TRIGGER AS $$
BEGIN
    -- ตั้งค่า Updated_At ให้เป็นเวลาปัจจุบัน
    NEW.updated_at := CURRENT_TIMESTAMP;
    
    -- ส่งแถวที่ได้รับการอัปเดตกลับ
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE OR REPLACE TRIGGER updated_product
AFTER UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_insert();