const apiAdapter = require("../apiAdapter");

const api = apiAdapter(process.env.URL_SERVICE_FOOD)

module.exports = async (req, res) => {
    console.log("Hallo");
    try {
        const food = await api.patch(`/food/${req.params.id}`, req.body)
        return res.json(food.data)
    } catch (err) {
        if (err.code === 'ECONNREFUSED') {
            return res.status(500).json({ status: 'error', message: 'Service unavailable' });
        }

        const { status, data } = err.response;
        return res.status(status).json(data);
    }
}