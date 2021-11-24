const apiAdapter = require("../apiAdapter");

const api = apiAdapter(process.env.URL_SERVICE_FOOD)

module.exports = async (req, res) => {
    try {
        const food = await api.delete(`/food/${req.params.id}`);
        return res.json(food.data)
    } catch (err) {
        if (err.code === 'ECONNREFUSED') {
            return next(new AppError("Service unavailable.", 500));
        }
        const { status, data } = err.response;
        return res.status(status).json(data);
    }
}