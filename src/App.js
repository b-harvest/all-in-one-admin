import {
  BrowserRouter as Router,
  Switch,
  Route,
  // Link
} from "react-router-dom";
import Main from '@page/main'
import AddressFileReader from './components/fileReader/fileReader.js';
import AppHeader from './components/header/header.js'
import BasicLayout from '@layout/Basic'

function App() {
  return (
    <Router>
      <Switch>
        <Route exact path={["/", "/main"]}>
          <BasicLayout>
            <Main />
          </BasicLayout>
        </Route>
        <Route exact path="/chainlink">
          <AppHeader />
          <AddressFileReader />
        </Route>
      </Switch>
    </Router>
  );
}

export default App;
