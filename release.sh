#!/bin/bash

# TODO Turn this script into an action and rpc call the github api to run
# the action with parameters from this script that will shrink to a two liner

set -e
set -x

# current git branch
branch=$(git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')

# current project name
projectName=$(git config --local remote.origin.url|sed -n 's#.*/\([^.]*\)\.git#\1#p')

# set master branch
masterBranch=$branch

# release version, set from command line
releaseVersion=""

# version to set right after the release, set from command line
postReleaseVersion=""

# v1.0.0, v1.7.8, etc..
versionLabel=""

# release branch name
releaseBranch=""

# release branch name
tagName=""

# file in which to update version number
versionFile="version.go"

init ()
{
		# master branch validation
	if [ $branch != "master" ]; then
		echo "needs to run on master branch"
		exit -1
	fi

	if [ "$releaseVersion" = "" ]; then
		echo "Enter the release version number without a leading v"
		read releaseVersion
	fi

	if [ "$postReleaseVersion" = "" ]; then
		echo "Enter version number to set for $projectName immediately after the release"
		read postReleaseVersion
	fi

	versionLabel="v$releaseVersion"
	releaseBranch=release-$releaseVersion
	tagName=$versionLabel
	versionFile="version.go"

	echo "release version: $releaseVersion"
	echo "version label: $versionLabel"
	echo "release branch: $releaseBranch"
	echo "tag name: $tagName"
	echo "version file: $versionFile"
}

prerelease ()
{
	echo "Started releasing $versionLabel for $projectName ..."
	# pull the latest version of the code from master
	git pull
 
	# find version number ("1.5.5" for example)
	# and replace it with newly specified version number
	sed -i.backup -E "s/[0-9.]+[0-9]+/$version/" $versionFile $versionFile

	# remove backup file created by sed command
	rm -f $versionFile.backup
	
	# adding the changes to git
	git add $versionFile

	# commit version number increment
	git commit -m "[Pre Release] Setting version to $versionLabel, Creating branch $releaseBranch"
}

postrelease ()
{
		# create the release branch from the -master branch
	git checkout -b $releaseBranch $masterBranch

	# push local releaseBranch to remote
	git push -u origin $releaseBranch

	echo "$projectName $versionLabel successfully released and labeled as $versionLabel"
	echo "Incrementing $masterBranch version for all builds after the release"

	# checkout to master branch
	git checkout $masterBranch

	# pull the latest version of the code from master
	git pull
	
	# find version number ("1.5.5" for example)
	# and replace it with newly specified version number
	sed -i.backup -E "s/[0-9.]+[0-9]+/$postReleaseVersion/" $versionFile $versionFile

	# remove backup file created by sed command
	rm -f $versionFile.backup
	
	# adding the changes to git
	git add $versionFile

	# Commit setting new master branch version	
	git commit -m "[Post Release] Setting version to $postReleaseVersion"
}

tag ()
{
	# create tag for new version from -master
	git tag -f $tagName
}

push ()
{
	# push tag to remote origin. This will actually trigger CI
	git push --tags origin
}

release ()
{
	local releaseVersion=$1
	local postReleaseVersion=$2

	init
	echo "[$1] -> [$2]"
	prerelease "$releaseVersion"
	tag "$versionLabel"
	push
	postrelease "$postReleaseVersion"
	push
}

release $1 $2