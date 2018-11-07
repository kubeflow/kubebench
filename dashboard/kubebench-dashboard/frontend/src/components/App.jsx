import React  from 'react';
import Header from './Header';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';

import GenerateFromYaml from './GenerateFromYaml';
import GenerateFromParameters from './GenerateFromParameters';
import Watch from './Watch';

import  { Route, Link } from 'react-router-dom';


const styles = {
    root: {
        width: '90%',
        margin: '0 auto',
        paddingTop: 20,
    }
};

const App = (props) => {
    const { classes } = props;
    return (
        <div className={classes.root}>
            <Header />

            <ul>
                <li>
                    <Link to="/">Generate from YAML</Link>
                </li>
                <li>
                    <Link to="/defaults">Generate from parameters</Link>
                </li>
                <li>
                    <Link to="/watch">Watch</Link>
                </li>
            </ul>

            <Route exact path="/" component={GenerateFromYaml} />
            <Route path="/defaults" component={GenerateFromParameters} />
            <Route path="/watch" component={Watch} />
        </div>
    )
};

App.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(App)
