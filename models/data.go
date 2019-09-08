package models

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func GetOrderInfo(sale_no string) ([]map[string]interface{}) {

	rows, err := DB.Query(`SELECT * FROM order_pos_infos WHERE transaction_order_id = ?`, sale_no)

	if err != nil {
		log.Fatalf("[AppSettingModel GetAppSettings Err]: %s\n", err)
	}

	output := make([]map[string]interface{}, 0)

	var order_id int
    var transaction_order_id string 
    var table_no string 
    var pax sql.NullInt64

	for rows.Next() {

        err := rows.Scan(&order_id, &transaction_order_id, &table_no, &pax)

        if err != nil {
            log.Fatalf("[AppUiImagesRowScanError]: %s\n", err)
        }

        content := map[string]interface{}{
            "order_id": order_id,
            "transaction_order_id": transaction_order_id,
            "table_no": table_no,
            "pax": pax,
        }

        output = append(output, content)
    }

	return output
}

func UpdateOrderStatus(order_id int, status int) {

	_, err := DB.Exec(`UPDATE orders SET status = ? WHERE id = ?`, status, order_id)

	if err != nil {
        log.Fatalf("[AppUiImagesRowScanError]: %s\n", err)
    }
}