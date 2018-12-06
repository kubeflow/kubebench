import React from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import { closeSnack, changeParameters, submitYamlParameters } from '../actions';

import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import Tooltip from '@material-ui/core/Tooltip';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import Button from "@material-ui/core/Button";


const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
    },
    editor: {
        margin: '0 auto',
    },
    submit: {
        textAlign: 'center',
    },
    button: {
        margin: theme.spacing.unit,
    },
    progress: {
        height: 10,
        margin: 10,
    },
    close: {
        padding: theme.spacing.unit / 2,
    },
    textField: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
        width: '100%'
    },
    help: {
        padding: theme.spacing.unit / 2,
        verticalAlign: "middle",
    },
    section: {
        padding: theme.spacing.unit,
    },
    parameter: {
        padding: theme.spacing.unit / 2,
    }
});


const GenerateFromParameters = (props) => {

    const onFieldChange = (name) => (event) => {
        return props.changeParameters(name, event.target.value);
    };

    const submitYamlParameters = () => {
        props.submitYamlParameters(props.parameters);
    };

    const { classes } = props;

    return (
        <div className={classes.root}>
            <h1>Generate from parameters</h1>
            <hr />
            {props.loading && <LinearProgress className={classes.progress}/>}
            {
                props.parameters.map((param, i) => {
                    if (param.description === "section") {
                        return (
                            <div key={i} className={classes.section}>
                                <Grid container>
                                    <Grid item xs={12} sm={12}>
                                        <Typography variant="h6">
                                            {param.name}
                                        </Typography>
                                        <hr />
                                    </Grid>
                                </Grid>
                            </div>
                        )
                    }
                    return (
                        <div key={i} className={classes.parameter}>
                            <Grid container>
                                <Grid item xs={12} sm={3}>
                                    <Typography>
                                        <Tooltip title={param.description}>
                                            <HelpOutlineIcon className={classes.help} color={"secondary"}/>
                                        </Tooltip>
                                        {param.name}
                                    </Typography>
                                </Grid>
                                <Grid item xs={12} sm={8}>
                                    <TextField
                                        className={classes.textField}
                                        value={param.value}
                                        onChange={onFieldChange(param.name)}
                                        />
                                </Grid>
                            </Grid>
                        </div>
                    )
                })
            }
            <div className={classes.submit}>
                <Button variant="contained" disabled={props.loading} color={"primary"} className={classes.button} onClick={submitYamlParameters}>
                    Deploy
                </Button>
            </div>
            <Snackbar
                anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'center',
                }}
                open={props.snackOpen}
                autoHideDuration={6000}
                onClose={props.closeSnack}
                message={<span id="message-id">{props.snackText}</span>}
                action={[
                    <IconButton
                        key="close"
                        aria-label="Close"
                        color="inherit"
                        className={classes.close}
                        onClick={props.closeSnack}
                    >
                        <CloseIcon />
                    </IconButton>
                ]}
            />
        </div>
    );
};

const mapStateToProps = (state) => {
    return {
        loading: state.loading,
        snackOpen: state.snackOpen,
        snackText: state.snackText,
        parameters: state.parameters,
    };
};

GenerateFromParameters.propTypes = {
    classes: PropTypes.object.isRequired,
};


export default connect(mapStateToProps, { closeSnack, changeParameters, submitYamlParameters })(withStyles(styles)(GenerateFromParameters));