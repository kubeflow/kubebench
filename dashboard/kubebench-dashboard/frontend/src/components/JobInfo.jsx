import React from 'react';
import Modal from '@material-ui/core/Modal';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Button from "@material-ui/core/Button";
import { deleteJob } from '../actions';
import {connect} from 'react-redux';

import CallSplit from '@material-ui/icons/CallSplit';
import FolderIcon from '@material-ui/icons/Folder'
import InsertChart from '@material-ui/icons/InsertChart'


const styles = theme => ({
    paper: {
        position: 'absolute',
        width: theme.spacing.unit * 50,
        backgroundColor: theme.palette.background.paper,
        boxShadow: theme.shadows[5],
        padding: theme.spacing.unit * 4,
        textAlign: 'center'
    },
    modal: {
        top: '40%',
        left: '40%',
    },
    button: {
        margin: theme.spacing.unit,
        width: '50%',
    },
    extraIcon: {
        marginRight: theme.spacing.unit,
    }
});


const JobInfo = (props) => {
    const { classes } = props;

    const deleteJob = () => {
        props.deleteJob(props.name);
    };

    return (
        <Modal
            aria-labelledby="simple-modal-title"
            aria-describedby="simple-modal-description"
            open={props.open}
            onClose={props.close}
            className={classes.modal}
        >
            <div className={classes.paper}>
                <Typography variant="h3" id="modal-title">
                    { props.name }
                </Typography>
                <Typography variant="subtitle1" id="simple-modal-description">
                    <a href={`/argo/workflows/default/${props.name}`} target={"_blank"}>
                        <Button variant="extendedFab" aria-label="Delete" className={classes.button} color={"secondary"}>
                            <CallSplit className={classes.extraIcon} />
                            Workflow
                        </Button>
                    </a>
                    <br />
                    <a href={`/grafana/d/eqSbm0Aik3/kubebench-monitoring`} target={"_blank"}>
                        <Button variant="extendedFab" className={classes.button} color={"secondary"}>
                            <InsertChart className={classes.extraIcon} />
                            Metrics
                        </Button>
                    </a>
                    <br />
                    <a href={`/`} target={"_blank"}>
                        <Button variant="extendedFab" className={classes.button} color={"secondary"}>
                            <FolderIcon className={classes.extraIcon} />
                            NFS
                        </Button>
                    </a>
                    <Button variant="extendedFab" className={classes.button} color={"primary"} onClick={deleteJob}>
                        Delete
                    </Button>
                </Typography>
            </div>
        </Modal>
    );
};

JobInfo.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default connect(null, { deleteJob })(withStyles(styles)(JobInfo));