func XenditWebhook(c echo.Context) error {
    token := c.Request().Header.Get("x-callback-token")

    if token != os.Getenv("XENDIT_CALLBACK_TOKEN") {
        return c.JSON(401, map[string]string{
            "message": "invalid token",
        })
    }

    payload := new(XenditPayload)

    if err := c.Bind(payload); err != nil {
        return err
    }

    // update transaksi
    err := transactionUsecase.MarkAsPaid(payload.ExternalID)

    if err != nil {
        return err
    }

    // publish rabbitmq
    rabbit.Publish("payment.paid", payload)

    return c.JSON(200, map[string]string{
        "message": "success",
    })
}