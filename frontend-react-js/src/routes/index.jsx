import { createBrowserRouter } from "react-router-dom";

import App from "../App";
import Error from "../pages/Error";
import Home from "../pages/Home";
import Login from "../pages/Login";
import ConfirmSignup from "../pages/ConfirmSignup";
import NewPass from "../pages/NewPass";
import ResetPass from "../pages/ResetPass";
import Pricing from "../pages/Pricing";
import Products from "../pages/Products";
import Admin from "../pages/Admin";
import Checkout from "../pages/Checkout";
import AddBalance from "../pages/AddBalance";
import RoutePrivate from "../components/PrivateRoute";
import ConfirmEmail from "../pages/ConfirmEmail";
import Profile from "../pages/Profile";
import ResetPassConfirm from "../pages/ResetPassConfirm";

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
        path: "confirm-signup",
        element: <ConfirmSignup />,
      },
      {
        path: "confirm-email",
        element: <ConfirmEmail />,
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
        path: "resetpass-confirm",
        element: <ResetPassConfirm />,
      },
      {
        path: "profile",
        element: (
          <RoutePrivate>
            <Profile />
          </RoutePrivate>
        ),
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
          <RoutePrivate>
            <Checkout />
          </RoutePrivate>
        ),
      },
      {
        path: "add-balance",
        element: (
          <RoutePrivate>
            <AddBalance />
          </RoutePrivate>
        ),
      },
    ],
  },
]);
