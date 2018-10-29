// @apiVersion 0.1
// @name io.ksonnet.pkg.kubebench-dashboard
// @description Kubebench dashboard installer
// @shortDescription Kubebench dashboard installer
// @param name string Name for the component
// @optionalParam image string docker.io/akado2009/kb-dashboard:latest Image for kubebench dashboard

local k = import "k.libsonnet";

local kubebenchDashboard = import "kubebench/kubebench-dashboard/kubebench-dashboard.libsonnet";

local kubebenchDashboardInstance = kubebenchDashboard.new(env, params);
kubebenchDashboardInstance.list(kubebenchDashboardInstance.all)