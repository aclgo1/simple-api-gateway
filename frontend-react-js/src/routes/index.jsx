import { createBrowserRouter } from "react-router-dom";

import App from "../App";
import Error from "../pages/Error";
import Home from "../pages/Home";
import Login from "../pages/Login";
import ConfirmSignup from "../pages/C";
import NewPass from "../pages/NewPass";
import ResetPass from "../pages/ResetPass";
import Pricing from "../pages/Pricing";
import Products from "../pages/Products";
import Admin from "../pages/Admin";
import Checkout from "../pages/Checkout";
import AddBalance from "../pages/AddBalance";
import RoutePrivate from "../components/PrivateRoute";

export const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    errorElement: <Error />,
    children: [
      {
        index: true,
        element: <Pricing />,
      },
      {
        path: "home",
        element: (
          <RoutePrivate>
            <Home />
          </RoutePrivate>
        ),
      },
      {
        path: "login",
        element: <Login />,
      },
      {
        path: "confirm",
        element: <ConfirmSignup />,
      },
      {
        path: "newpass",
        element: <NewPass />,
      },
      {
        path: "resetpass",
        element: <ResetPass />,
      },
      {
        path: "products",
        element: (
          <RoutePrivate>
            <Products />
          </RoutePrivate>
        ),
      },
      {
        path: "admin",
        element: (
          <RoutePrivate>
            <Admin />
          </RoutePrivate>
        ),
      },
      {
        path: "checkout",
        element: (
          // <RoutePrivate>
          <Checkout />
          // </RoutePrivate>
        ),
      },
      {
        path: "balance",
        element: (
          <RoutePrivate>
            <AddBalance />
          </RoutePrivate>
        ),
      },
    ],
  },
]);
