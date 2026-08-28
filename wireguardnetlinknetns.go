//go:build wireguard_netlink_netns

package main

import (
	"github.com/arnika-project/arnika/config"
	"github.com/arnika-project/arnika/repositories"
	"github.com/arnika-project/arnika/services"
)

func getKeyWriterService(cfg *config.Config) (*services.KeyWriterService, error) {
	wireguardRepo, err := repositories.NewWireguardNetlinkNetnsRepository(cfg.WireGuardInterface, cfg.WireguardPeerPublicKey, cfg.WireGuardNetnsPath)
	if err != nil {
		return nil, err
	}
	return services.NewKeyWriterService(wireguardRepo), nil
}
