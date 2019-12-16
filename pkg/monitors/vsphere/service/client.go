package service

import (
	"context"
	"crypto/tls"
	"net/url"

	"github.com/signalfx/signalfx-agent/pkg/monitors/vsphere/model"
	"github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// LogIn logs into vCenter and returns a logged-in Client or an error
func LogIn(ctx context.Context, conf *model.Config) (*govmomi.Client, error) {
	myUrl, err := url.Parse("https://" + conf.Host + "/sdk")
	if err != nil {
		return nil, err
	}
	myUrl.User = url.UserPassword(conf.Username, conf.Password)

	Log.WithFields(logrus.Fields{
		"ip":   conf.Host,
		"user": conf.Username,
	}).Info("Connecting to vsphereInfo")

	client, err := newGovmomiClient(ctx, myUrl, conf)
	if err != nil {
		return nil, err
	}

	err = client.Login(ctx, myUrl.User)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newGovmomiClient(ctx context.Context, myUrl *url.URL, conf *model.Config) (*govmomi.Client, error) {
	vimClient, err := newVimClient(ctx, myUrl, conf)
	if err != nil {
		return nil, err
	}
	return &govmomi.Client{
		Client:         vimClient,
		SessionManager: session.NewManager(vimClient),
	}, nil
}

func newVimClient(ctx context.Context, myUrl *url.URL, conf *model.Config) (*vim25.Client, error) {
	soapClient := soap.NewClient(myUrl, conf.InsecureSkipVerify)
	if conf.TLSCACertPath != "" {
		Log.Info("Attempting to load TLSCACertPath from ", conf.TLSCACertPath)
		err := soapClient.SetRootCAs(conf.TLSCACertPath)
		if err != nil {
			return nil, err
		}
	} else {
		Log.Info("No tlsCACertPath provided. Not setting root CA.")
	}
	if conf.TLSClientCertificatePath != "" && conf.TLSClientKeyPath != "" {
		Log.Infof(
			"Attempting to load client certificate from TLSClientCertificatePath(%s) and TLSClientKeyPath(%s)",
			conf.TLSClientCertificatePath,
			conf.TLSClientKeyPath,
		)
		cert, err := tls.LoadX509KeyPair(conf.TLSClientCertificatePath, conf.TLSClientKeyPath)
		if err != nil {
			return nil, err
		}
		soapClient.SetCertificate(cert)
	} else {
		Log.Info("Either or both of tlsClientCertificatePath or tlsClientKeyPath not provided. Not setting client certificate.")
	}
	return vim25.NewClient(ctx, soapClient)
}
