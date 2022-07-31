import 'package:flutter/material.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';

enum FriendStatus {
  free,
  busy,
  invisible,
}

class FriendCard extends StatelessWidget {
  final String name;
  final String pic;
  final FriendStatus status;
  // curEventName is the name of the current event the friend is attending
  final String curEventName;
  final VoidCallback showOverflowMenu;
  final GestureTapCallback? onTap;

  const FriendCard({
    Key? key,
    required this.name,
    this.pic = 'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
    required this.status,
    required this.showOverflowMenu,
    this.curEventName = '',
    this.onTap,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: SchejConstants.listTileDecoration,
      child: ListTile(
        onTap: onTap,
        dense: true,
        leading: CircleAvatar(
          backgroundImage: NetworkImage(pic),
        ),
        title: Text(name, style: SchejFonts.subtitle),
        subtitle: _buildStatusText(),
        trailing: IconButton(
          icon: const Icon(MdiIcons.dotsVertical),
          onPressed: showOverflowMenu,
        ),
      ),
    );
  }

  Widget _buildStatusText() {
    if (status == FriendStatus.free) {
      return Text(
        'Currently free',
        style: SchejFonts.bodyMedium.copyWith(color: SchejColors.darkGreen),
      );
    } else if (status == FriendStatus.invisible) {
      return Text(
        'Currently invisible',
        style: SchejFonts.body.copyWith(color: SchejColors.darkGray),
      );
    } else if (status == FriendStatus.busy) {
      return Text(
        'Currently at $curEventName',
        style: SchejFonts.body.copyWith(color: SchejColors.pureBlack),
      );
    }

    return const Text('');
  }
}
