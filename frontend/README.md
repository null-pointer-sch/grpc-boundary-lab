# gRPC Boundary Lab - Frontend Dashboard

This directory contains the React/Vite visualization dashboard for the gRPC Boundary Lab.

It provides a live, interactive UI to trigger the latency `ping` benchmarking endpoints and visualize `make sweep` statistics, allowing users to quickly toggle between observing native **gRPC** proxy performance and simulated **REST** proxy performance.

## Tech Stack
- React 18
- Vite
- Tailwind CSS
- Recharts
- Axios

## Development
To run the frontend independently:
```bash
npm install
npm run dev
```
*(Requires the `backend` and `gateway` to be running on their respective ports for the application to function).*

## Docker 
The recommended execution method is via Docker Compose from the root directory (`make all` or `docker compose up --build -d`), which bundles this React app into an `Nginx` container serving on port `80`.
