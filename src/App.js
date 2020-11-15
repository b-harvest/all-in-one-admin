import {
  BrowserRouter as Router,
  Switch,
  Route,
  // Link
} from "react-router-dom";
import Main from '@page/main'
import BasicLayout from '@layout/Basic'

function App() {
  return (
    <Router>
      <Switch>
        <Route exact path="/">
          <BasicLayout>
            <Main />
          </BasicLayout>
        </Route>
      </Switch>
    </Router>
  );
}

export default App;
