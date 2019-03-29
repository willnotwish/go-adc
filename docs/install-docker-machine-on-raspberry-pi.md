### How to set up Docker Machine on raspberry pi
1. Download Raspbian
2. Use Etcher to flash image onto SD card
3. Eject & reinsert SD card
4. Locate new volume boot using Finder
5. Create an empty file called ssh in the root of the boot volume using touch.
6. Eject SD card
7. Insert into the pi, ensure it’s connected to the LAN and switch it on
8. Locate the pi on the network with something like
sudo nmap -sP 192.168.0.0/24 | awk '/^Nmap/{ip=$NF}/B8:27:EB/{print ip}'
9. ssh to pi as pi/raspberry
10. Add .ssh directory
mkdir ~/.ssh
11. Copying your public key from laptop to pi
cat ~/.ssh/id_rsa.pub | ssh pi@192.168.0.50 'cat >> .ssh/authorized_keys'
12. Log in to pi without a password
ssh pi@192.168.0.50
13. Edit /etc/os-release to change ID to debian instead of raspbian.
14. Reboot pi
15. ssh back in
16. Install docker
curl -sSL https://get.docker.com | sh
17. Add user pi to group docker
sudo usermod -aG docker pi
18. Log out
19. Check docker is running on pi
ssh pi@192.168.0.50 docker version
20. Run docker hello-world
ssh pi@192.168.0.50 docker run hello-world
21. Provision docker-machine[a]
docker-machine create --driver generic --generic-ip-address=192.168.0.50 --generic-ssh-key ~/.ssh/id_rsa --generic-ssh-user pi --engine-storage-driver overlay2 pi3b
22. Check docker-machine
docker-machine ls

See https://gist.github.com/calebbrewer/c41cab61216d8845b59fcc51f36343a7

[a]Do *not* let docker-machine provision the pi or it will fail. Use the overlay file system and not AUFS.