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
  final FriendStatus status;
  // curEventName is the name of the current event the friend is attending
  final String curEventName;
  final VoidCallback showOverflowMenu;

  const FriendCard({
    Key? key,
    required this.name,
    required this.status,
    required this.showOverflowMenu,
    this.curEventName = '',
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: SchejColors.white,
        borderRadius: SchejConstants.borderRadius,
        border: Border.all(
          color: SchejColors.offWhite,
          width: 1,
        ),
        boxShadow: [
          BoxShadow(
            color: SchejColors.pureBlack.withOpacity(0.1),
            spreadRadius: 0,
            blurRadius: 4,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: ListTile(
        dense: true,
        leading: const CircleAvatar(
          backgroundImage: NetworkImage(
            'https://pbs.twimg.com/media/D8dDZukXUAAXLdY.jpg',
          ),
        ),
        title: Text(
          name,
          style: SchejFonts.subtitle.copyWith(color: SchejColors.pureBlack),
        ),
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
