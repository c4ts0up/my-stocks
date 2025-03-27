# My-Stocks

### By [c4ts0up](http://c4ts0up.github.io)

---

## 📝 Description
A web app to retrieve, analyze and present stock 
information. Make sure to check out the [Wiki](https://github.com/c4ts0up/my-stocks/wiki) for more information.

## 🚀 Features
- Tracks live stock prices
- Displays stock information, including recent ratings 
from numerous rating agencies
- Analyzes stock ratings to offer a summarized 
recommendation

## 🛠️ Tech Stack
**Frontend:** Vue 3, TypeScript, TailwindCSS  
**Backend:** Go (Gin), GORM, CockroachDB  
**Deployment:** Docker, Docker Compose  

## ⚙️ Local Deployment
1. Clone the repository
2. Set up the `local.env` file in the root of the project
3. Modify `./frontend/.env` according to the local deployment
4. `docker compose up --build`

The backend will be deployed on ``http://localhost:8080`` while
the frontend will be deployed on ``http://localhost:5173``

⚠️ **NOTE**: This deployment is not production safe.

### `local.env`
This file is necessary to deploy locally. It must be located in the root of the project directory. It requires the following variables:
- `DATABASE_URL`
- `FETCH_DELAY_S`
- `ANALYSIS_DELAY_S`
- `INFO_API_URL`
- `INFO_API_TOKEN`
- `RATINGS_API_URL`
- `RATINGS_API_TOKEN`
- `FRONTEND_URL`
- `STOCKS_API_URL`
