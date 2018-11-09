import React  from 'react';
import Header from './Header';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';

import GenerateFromYaml from './GenerateFromYaml';
import GenerateFromParameters from './GenerateFromParameters';
import Watch from './Watch';

import  { Route, Link } from 'react-router-dom';
import { Menu, Icon } from 'antd';
import "antd/dist/antd.css";


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
            <Menu
                // selectedKeys={[this.state.current]}
                mode="horizontal"
            >
                <Menu.Item key="yaml">
                    <Link to="/">
                        <Icon type="file-text" />
                        Generate from YAML
                    </Link>
                </Menu.Item>
                <Menu.Item key="param">
                    <Link to="/defaults">
                        <Icon type="ordered-list" />
                        Generate from parameters
                    </Link>
                </Menu.Item>
                <Menu.Item key="watch">
                    <Link to="/monitor">
                        <Icon type="eye" />
                        Monitor
                    </Link>
                </Menu.Item>
            </Menu>

            <Route exact path="/" component={GenerateFromYaml} />
            <Route path="/defaults/" component={GenerateFromParameters} />
            <Route path="/monitor/" component={Watch} />
        </div>
    )
};

App.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(App)
