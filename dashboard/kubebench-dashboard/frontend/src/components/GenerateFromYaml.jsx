import React from 'react';
import {connect} from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';

import { changeYaml, submitWholeYaml, closeSnack } from '../actions';

import brace from 'brace';
import 'brace/mode/javascript';
import 'brace/theme/tomorrow';
import AceEditor from 'react-ace';
import Button from '@material-ui/core/Button';

import LinearProgress from '@material-ui/core/LinearProgress';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';


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
});


const GenerateFromYaml = (props) => {

    const onYamlChange = (value) => {
        props.changeYaml(value);
    };

    const submitWholeYaml = () => {
        props.submitWholeYaml(props.yaml);
    };

    const { classes } = props;
    return (
        <div className={classes.root}>
            <h1>Generate</h1>
            <hr />
            {props.loading && <LinearProgress className={classes.progress}/>}
            <div className={classes.editor}>
                <AceEditor
                    mode="text"
                    theme="tomorrow"
                    value={props.yaml}
                    onChange={onYamlChange}
                    name="UNIQUE_ID_OF_DIV"
                    editorProps={{$blockScrolling: true}}
                    tabSize={2}
                    enableLiveAutocompletion={true}
                    fontSize={14}
                    width={'100%'}
                    height={700}
                />
            </div>
            <div className={classes.submit}>
                <Button variant="contained" disabled={props.loading} color={"primary"} className={classes.button} onClick={submitWholeYaml}>
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
                    </IconButton>,
                ]}
            />
        </div>
    );
};

const mapStateToProps = (state) => {
    return {
        yaml: state.yaml,
        loading: state.loading,
        snackOpen: state.snackOpen,
        snackText: state.snackText,
    };
};

GenerateFromYaml.propTypes = {
    classes: PropTypes.object.isRequired,
};


export default connect(mapStateToProps, { changeYaml, submitWholeYaml, closeSnack })(withStyles(styles)(GenerateFromYaml));