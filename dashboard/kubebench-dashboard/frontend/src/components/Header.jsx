import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';


const styles = {
    root: {
        flexGrow: 1,
    },
    grow: {
        flexGrow: 1,
    },
};

const Header = (props) => {
    const { classes } = props;

    return (
        <div className={classes.root}>
            <AppBar position={"static"} color={"primary"}>
                <Toolbar>
                    <Typography variant={"display1"} color={"inherit"} className={classes.grow}>
                        <strong>Kubebench Dashboard</strong>
                    </Typography>
                </Toolbar>
            </AppBar>
        </div>
    );
};

Header.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Header);